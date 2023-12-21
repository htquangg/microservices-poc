package mysql

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"
)

type customer struct {
	ID        string    `xorm:"pk notnull"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Name      string    `xorm:"varchar(64) notnull"`
	Phone     string    `xorm:"varchar(16) notnull"`
	Email     string    `xorm:"varchar(64) null"`
}

func (customer) TableName() string {
	return "customers"
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	db database.DB
}

func NewCustomerRepository(db database.DB) CustomerRepository {
	return CustomerRepository{
		db: db,
	}
}

func (r CustomerRepository) Save(ctx context.Context, c *domain.Customer) error {
	return r.db.Insert(ctx, &customer{
		ID:    c.ID(),
		Name:  c.Name(),
		Phone: c.Phone(),
		Email: c.Email(),
	})
}

func (CustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	panic("unimplemented")
}

func (CustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	panic("unimplemented")
}
