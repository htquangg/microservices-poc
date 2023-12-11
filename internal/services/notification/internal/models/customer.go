package models

import "github.com/htquangg/microservices-poc/internal/services/notification/internal/constants"

type Customer struct {
	ID    string
	Name  string
	Email string
	Phone string
}

func (Customer) TableName() string {
	return constants.CustomerTableName
}
