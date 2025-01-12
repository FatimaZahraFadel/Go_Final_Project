package book

import (
	"context"
	"time"

	"um6p.ma/final_project/internal/author"
)

type Book struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Author      author.Author `json:"author"`
	Genre       string        `json:"genres"`
	PublishedAt time.Time     `json:"published_at"`
	Price       float64       `json:"price"`
	Stock       int           `json:"stock"`
}

type BookStore interface {
	CreateBook(ctx context.Context, book Book) (Book, error)
	GetBook(ctx context.Context, id int) (Book, error)
	UpdateBook(ctx context.Context, id int, book Book) (Book, error)
	DeleteBook(ctx context.Context, id int) error
	SearchBooks(ctx context.Context, criteria SearchCriteria) ([]Book, error)
	GetAllBooks(ctx context.Context) ([]Book, error)
}

type SearchCriteria struct {
	Title       string
	Author      author.Author
	Genre       string
	PublishedAt time.Time
	PriceRange  [2]float64
}
