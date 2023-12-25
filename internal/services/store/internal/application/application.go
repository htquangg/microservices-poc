package application

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/logger"
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

func New(
	storeESRepo domain.StoreESRepository,
	productESRepo domain.ProductESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) *Application {
	return &Application{
		Commands: Commands{
			CreateStoreHandler:  commands.NewCreateStoreHandler(storeESRepo, publisher, log),
			RebrandStoreHandler: commands.NewRebrandStoreHandler(storeESRepo, publisher, log),
			AddProductHandler:   commands.NewAddProductHandler(productESRepo, publisher, log),
		},
	}
}
