package system

import (
	"github.com/htquangg/microservices-poc/internal/services/payment/config"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/discovery"
	"github.com/htquangg/microservices-poc/pkg/logger"
	"github.com/htquangg/microservices-poc/pkg/waiter"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type Service interface {
	Config() *config.Config
	DB() database.DB
	Router() *mux.Router
	RPC() *grpc.Server
	Discovery() discovery.Registry
	Logger() logger.Logger
	Waiter() waiter.Waiter
}
