package mysql

import (
	"context"
	"fmt"
	"time"

	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const CustomerTable = "customers"

type customer struct {
	ID        string    `xorm:"pk notnull"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Name      string    `xorm:"varchar(64) notnull"`
	Phone     string    `xorm:"varchar(16) notnull"`
	Email     string    `xorm:"varchar(64) null"`
}

func (customer) TableName() string {
	return CustomerTable
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

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	query := r.table("SELECT name, phone, email from %s WHERE id = ? LIMIT 1")

	results, err := r.db.Engine(ctx).Query(query, customerID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, mysql_internal.ErrRecordNotFound
	}

	return domain.NewCustomer(
		customerID,
		string(results[0]["name"]),
		string(results[0]["phone"]),
		string(results[0]["email"]),
	)
}

func (CustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	panic("unimplemented")
}

func (CustomerRepository) table(query string) string {
	return fmt.Sprintf(query, CustomerTable)
}
