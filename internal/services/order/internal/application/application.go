package application

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		commands.CreateOrderHandler
		commands.CancelOrderHandler
	}

	Queries struct{}
)

func New(
	orderESRepo domain.OrderESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) *Application {
	return &Application{
		Commands: Commands{
			CreateOrderHandler: commands.NewCreateOrderHandler(orderESRepo, publisher, log),
			CancelOrderHandler: commands.NewCancelOrderHandler(orderESRepo, publisher, log),
		},
	}
}
