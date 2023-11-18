package system

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/htquangg/microservices-poc/internal/services/customer/config"
	"github.com/htquangg/microservices-poc/pkg/constants"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/discovery"
	"github.com/htquangg/microservices-poc/pkg/discovery/consul"
	"github.com/htquangg/microservices-poc/pkg/logger"
	"github.com/htquangg/microservices-poc/pkg/rpc"
	"github.com/htquangg/microservices-poc/pkg/uid"
	"github.com/htquangg/microservices-poc/pkg/waiter"
	"github.com/htquangg/microservices-poc/pkg/web"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type System struct {
	cfg *Config

	db *database.DB

	router    *mux.Router
	rpc       *grpc.Server
	discovery discovery.Registry

	sf     *uid.Sonyflake
	logger logger.Logger
	waiter waiter.Waiter

	isRunningHTTP bool
	isRunningRPC  bool
}

func New(cfg *config.Config) (*System, error) {
	s := &System{cfg: &Config{
		Config: cfg,
	}}

	s.initWaiter()
	s.initLogger()
	s.initSonyflake()

	if err := s.initDB(); err != nil {
		return nil, err
	}

	s.initRouter()
	s.initRPC()

	if err := s.initDiscovery(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *System) Config() *config.Config {
	return s.cfg.Config
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) initLogger() {
	s.logger = logger.NewZapLogger(&logger.LogConfig{
		Environment: s.cfg.Environment,
	})
}

func (s *System) Logger() logger.Logger {
	return s.logger
}

func (s *System) initSonyflake() {
	s.sf = uid.NewSonyflake()
}

func (s *System) Sonyflake() *uid.Sonyflake {
	return s.sf
}

func (s *System) initDB() (err error) {
	s.db, err = database.New(s.waiter.Context(), s.logger, s.cfg.Mysql)
	if err != nil {
		return err
	}
	defer s.logger.Infof("%s server is connected to database: tcp(%v:%v)/%v",
		s.cfg.Name,
		s.cfg.Mysql.Host,
		s.cfg.Mysql.Port,
		s.cfg.Mysql.Schema,
	)
	return nil
}

func (s *System) DB() *database.DB {
	return s.db
}

func (s *System) initRouter() {
	s.router = mux.NewRouter()
	s.router.Use(handlers.RecoveryHandler(
		handlers.RecoveryLogger(s.logger),
		handlers.PrintRecoveryStack(true),
	))
	s.router.Use(handlers.CompressHandler)
	s.router.HandleFunc("/healthz", web.HealthCheck).
		Methods(http.MethodGet, http.MethodHead)
}

func (s *System) Router() *mux.Router {
	return s.router
}

func (s *System) initRPC() {
	s.rpc = grpc.NewServer(
		grpc.UnaryInterceptor(kitgrpc.Interceptor),
	)

	if s.cfg.IsDevelopment() {
		reflection.Register(s.rpc)
	}

	grpc_health_v1.RegisterHealthServer(s.rpc, &rpc.HealthImpl{})
}

func (s *System) RPC() *grpc.Server {
	return s.rpc
}

func (s *System) initDiscovery() (err error) {
	s.discovery, err = consul.New(s.cfg.Consul, s.logger)
	return err
}

func (s *System) Discovery() discovery.Registry {
	return s.discovery
}

func (s *System) WaitForWeb(ctx context.Context) error {
	if s.isRunningHTTP {
		return fmt.Errorf("%s web server is already running", s.cfg.Name)
	}
	s.isRunningHTTP = true

	webServer := &http.Server{
		Addr:              s.cfg.Web.Address(),
		Handler:           s.router,
		ReadTimeout:       25 * time.Second,
		WriteTimeout:      40 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       5 * time.Minute,
		MaxHeaderBytes:    1 << 18, // 0.25 MB (262144 bytes)
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		s.logger.Infof(
			"%s web server is listening on port: %d",
			s.cfg.Name,
			s.cfg.Web.Port,
		)
		defer s.logger.Infof("%s web server shutdown", s.cfg.Name)
		if err := webServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.Infof("%s web server to be shutdown", s.cfg.Name)
		ctx, cancel := context.WithTimeout(
			context.Background(),
			constants.WaitShutdownDuration,
		)
		defer cancel()
		err := webServer.Shutdown(ctx)
		return err
	})

	return group.Wait()
}

func (s *System) WaitForRPC(ctx context.Context) error {
	if s.isRunningRPC {
		return fmt.Errorf("%s rpc server is already running", s.cfg.Name)
	}
	s.isRunningRPC = true

	listener, err := net.Listen("tcp", s.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		s.logger.Infof(
			"%s rpc server is listening on port: %d",
			s.cfg.Name,
			s.cfg.Rpc.Port,
		)
		defer s.logger.Infof("%s rpc server shutdown", s.cfg.Name)
		if err := s.RPC().Serve(listener); err != nil &&
			err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.Infof("%s rpc server to be shutdown", s.cfg.Name)
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(constants.WaitShutdownDuration)
		select {
		case <-timeout.C:
			// Force it to stop
			s.RPC().Stop()
			return fmt.Errorf(
				"%s rpc server failed to stop gracefully",
				s.cfg.Name,
			)
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (s *System) WaitForWebDiscover(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := s.discovery.RegisterHTTP(
			s.webName(),
			s.webID(),
			"/healthz",
			s.cfg.Web.Host,
			s.cfg.Web.Port,
			nil,
		)
		return err
	})

	group.Go(func() error {
		<-gCtx.Done()
		err := s.discovery.Deregister(s.webID())
		return err
	})

	return group.Wait()
}

func (s *System) WaitForRPCDiscover(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := s.discovery.RegisterRPC(
			s.rpcName(),
			s.rpcID(),
			"Health",
			s.cfg.Rpc.Host,
			s.cfg.Rpc.Port,
			nil,
		)
		return err
	})

	group.Go(func() error {
		<-gCtx.Done()
		err := s.discovery.Deregister(s.rpcID())
		return err
	})

	return group.Wait()
}
