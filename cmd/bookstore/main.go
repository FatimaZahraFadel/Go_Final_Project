package main

import (
	"log"
	"net/http"

	"um6p.ma/final_project/internal/author"
	"um6p.ma/final_project/internal/book"
	"um6p.ma/final_project/internal/customer"
	"um6p.ma/final_project/internal/order"
	"um6p.ma/final_project/internal/sales"
)

func main() {
	http.HandleFunc("/authors", author.AuthorsHandler)
	http.HandleFunc("/authors/", author.AuthorHandler)
	http.HandleFunc("/books", book.BooksHandler)
	http.HandleFunc("/books/", book.BookHandler)
	http.HandleFunc("/customers", customer.CustomersHandler)
	http.HandleFunc("/customers/", customer.CustomerHandler)

	http.HandleFunc("/orders", order.OrdersHandler)
	http.HandleFunc("/orders/", order.OrderHandler)

	http.HandleFunc("/sales/report", sales.SalesReportHandler)

	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

}
