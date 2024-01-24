package saga

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/sec"
	"github.com/htquangg/microservices-poc/internal/services/cosec/models"
	"github.com/htquangg/microservices-poc/internal/services/customer/customerpb"
	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
)

const (
	CreateOrderSagaName     = "cosec.CreateOrder"
	CreateOrderReplyChannel = "mall.cosec.replies.CreateOrder"
)

type createOrderSaga struct {
	sec.Saga[*models.CreateOrderData]
}

func NewCreateOrderSaga() sec.Saga[*models.CreateOrderData] {
	saga := &createOrderSaga{
		Saga: sec.NewSaga[*models.CreateOrderData](CreateOrderSagaName, CreateOrderReplyChannel),
	}

	// 0. -RejectOrder
	saga.AddStep().
		Compensation(saga.rejectOrder)

	// 1. -AuthorizeCustomer
	saga.AddStep().
		Action(saga.authorizeCustomer)

	// 2. -ConfirmPayment
	saga.AddStep().
		Action(saga.confirmPayment).
		Compensation(saga.refundPayment)

	// 3. -ApproveOrder
	saga.AddStep().
		Action(saga.approveOrder)

	return saga
}

func (s *createOrderSaga) rejectOrder(ctx context.Context, data *models.CreateOrderData) (string, ddd.Command, error) {
	return orderpb.CommandChannel,
		ddd.NewCommand(orderpb.RejectOrderCommand, &orderpb.RejectOrder{
			Id: data.OrderID,
		}),
		nil
}

func (s *createOrderSaga) authorizeCustomer(
	ctx context.Context,
	data *models.CreateOrderData,
) (string, ddd.Command, error) {
	return customerpb.CommandChannel,
		ddd.NewCommand(customerpb.AuthorizeCustomerCommand, &customerpb.AuthorizeCustomer{
			Id: data.CustomerID,
		}),
		nil
}

func (s *createOrderSaga) confirmPayment(
	ctx context.Context,
	data *models.CreateOrderData,
) (string, ddd.Command, error) {
	panic("unimplemented")
}

func (s *createOrderSaga) refundPayment(
	ctx context.Context,
	data *models.CreateOrderData,
) (string, ddd.Command, error) {
	panic("unimplemented")
}

func (s *createOrderSaga) approveOrder(ctx context.Context, data *models.CreateOrderData) (string, ddd.Command, error) {
	return orderpb.CommandChannel, ddd.NewCommand(orderpb.ApproveOrderCommand, &orderpb.ApproveOrder{
		Id:         data.OrderID,
		ShoppingId: "",
	}), nil
}
