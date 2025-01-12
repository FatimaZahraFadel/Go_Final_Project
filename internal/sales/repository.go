package sales

import (
	"context"
	"time"

	"um6p.ma/final_project/internal/book"
)

type SalesReport struct {
	Timestamp       time.Time   `json:"timestamp"`
	TotalRevenue    float64     `json:"total_revenue"`
	TotalOrders     int         `json:"total_orders"`
	TopSellingBooks []BookSales `json:"top_selling_books"`
}

type BookSales struct {
	Book     book.Book `json:"book"`
	Quantity int       `json:"quantity_sold"`
}

type SalesStore interface {
	RecordSale(ctx context.Context, sale BookSales) error
	generateSalesReport(ctx context.Context, start, end time.Time) (SalesReport, error)
}
