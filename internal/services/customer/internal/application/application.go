package application

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application/command"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		RegisterCustomerHandler command.RegisterCustomerHandler
	}

	Queries struct{}
)

func New(customerRepo domain.CustomerRepository, publisher ddd.EventPublisher[ddd.AggregateEvent], log logger.Logger) *Application {
	return &Application{
		Commands: Commands{
			RegisterCustomerHandler: command.NewRegisterCustomerHandler(customerRepo, publisher, log),
		},
	}
}
