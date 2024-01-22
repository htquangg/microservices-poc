package application

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		commands.RegisterCustomerHandler
		commands.AuthorizeCustomerHandler
	}

	Queries struct{}
)

func New(
	customerRepo domain.CustomerRepository,
	publisher ddd.EventPublisher[ddd.AggregateEvent],
	log logger.Logger,
) *Application {
	return &Application{
		Commands: Commands{
			RegisterCustomerHandler:  commands.NewRegisterCustomerHandler(customerRepo, publisher, log),
			AuthorizeCustomerHandler: commands.NewAuthorizeCustomerHandler(customerRepo, publisher, log),
		},
	}
}
