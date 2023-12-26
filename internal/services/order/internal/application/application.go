package application

import (
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application/commands"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		commands.CreateStoreHandler
		commands.RebrandStoreHandler
		commands.AddProductHandler
	}

	Queries struct{}
)

func New() *Application {
	return &Application{
		Commands: Commands{},
	}
}
