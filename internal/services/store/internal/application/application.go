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
		commands.AddProductHandler
	}

	Queries struct{}
)

func New(
	productESRepo domain.ProductESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) *Application {
	return &Application{
		Commands: Commands{
			AddProductHandler: commands.NewAddProductHandler(productESRepo, publisher, log),
		},
	}
}
