package cmd

import (
	"fmt"
	"os"

	"github.com/htquangg/microservices-poc/internal/services/cosec/config"
	"github.com/htquangg/microservices-poc/internal/services/cosec/internal/system"
)

func Execute() {
	if err := run(); err != nil {
		fmt.Printf("store service exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg *config.Config

	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	s, err := system.New(cfg)
	if err != nil {
		return err
	}

	if err = startUp(s.Waiter().Context(), s); err != nil {
		return err
	}

	s.Logger().Infof("started %s service", cfg.Name)
	defer s.Logger().Infof("stopped %s service", cfg.Name)

	s.Waiter().Add(s.WaitForWeb, s.WaitForRPC, s.WaitForWebDiscover, s.WaitForRPCDiscover)

	return s.Waiter().Wait()
}
