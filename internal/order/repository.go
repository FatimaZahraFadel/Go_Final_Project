package order

import (
	"context"
	"time"

	"um6p.ma/final_project/internal/book"
	"um6p.ma/final_project/internal/customer"
)

type OrderItem struct {
	Book     book.Book `json:"book"`
	Quantity int       `json:"quantity"`
}

type Order struct {
	ID         int               `json:"id"`
	Customer   customer.Customer `json:"customer"`
	Items      []OrderItem       `json:"items"`
	TotalPrice float64           `json:"total_price"`
	CreatedAt  time.Time         `json:"created_at"`
	Status     string            `json:"status"`
}

type OrderStore interface {
	Create(ctx context.Context, order Order) (Order, error)
	GetByID(ctx context.Context, id int) (Order, error)
	Update(ctx context.Context, id int, order Order) (Order, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]Order, error)

	GetOrdersInTimeRange(ctx context.Context, start, end time.Time) ([]Order, error)
}
