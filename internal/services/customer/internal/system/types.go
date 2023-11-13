package system

import (
	"github.com/gorilla/mux"
	"github.com/htquangg/microservices-poc/internal/services/customer/config"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/discovery"
	"github.com/htquangg/microservices-poc/pkg/logger"
	"github.com/htquangg/microservices-poc/pkg/uid"
	"github.com/htquangg/microservices-poc/pkg/waiter"

	"google.golang.org/grpc"
)

type Service interface {
	Config() *config.Config
	DB() *database.DB
	Router() *mux.Router
	RPC() *grpc.Server
	Discovery() discovery.Registry
	Sonyflake() *uid.Sonyflake
	Logger() logger.Logger
	Waiter() waiter.Waiter
}
