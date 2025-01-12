package customer

import (
	"context"
	"time"
)

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   Address   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomerStore interface {
	GetCustomerByID(ctx context.Context, id int) (Customer, error)
	CreateCustomer(ctx context.Context, customer *Customer) (Customer, error)
	UpdateCustomer(ctx context.Context, id int, customer *Customer) (Customer, error) // Note: `*Customer`
	DeleteCustomer(ctx context.Context, id int) error
	GetAllCustomers(ctx context.Context) ([]Customer, error)
}
