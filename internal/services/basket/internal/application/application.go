package application

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		commands.StartBasketHandler
		commands.CancelBasketHandler
		commands.CheckoutBasketHandler
		commands.AddItemHandler
		commands.RemoveItemHandler
	}

	Queries struct{}
)

func New(
	basketESRepo domain.BasketESRepository,
	storeRepo domain.StoreRepository,
	productRepo domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) *Application {
	return &Application{
		Commands: Commands{
			StartBasketHandler:    commands.NewStartBasketHandler(basketESRepo, publisher, log),
			CancelBasketHandler:   commands.NewCancelBasketHandler(basketESRepo, publisher, log),
			CheckoutBasketHandler: commands.NewCheckoutBasketHandler(basketESRepo, publisher, log),
			AddItemHandler:        commands.NewAddItemHandler(basketESRepo, storeRepo, productRepo, publisher, log),
			RemoveItemHandler:     commands.NewRemoveItemHandler(basketESRepo, productRepo, publisher, log),
		},
	}
}
