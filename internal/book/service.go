package book

import (
	"context"
	"fmt"
	"log"
	"sync"

	"um6p.ma/final_project/internal/author"
)

type InMemoryBookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
}

func (store *InMemoryBookStore) CreateBook(ctx context.Context, book Book) (Book, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, existingbook := range store.books {
		if existingbook.Title == book.Title {
			log.Printf("book with title %s already exists", book.Title)
			return Book{}, fmt.Errorf("book with title %s already exists", book.Title)
		}
	}

	book.ID = store.nextID
	store.books[book.ID] = book
	store.nextID++

	return book, nil
}

func (store *InMemoryBookStore) GetBook(ctx context.Context, id int) (Book, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	book, found := store.books[id]
	if !found {
		log.Printf("book with ID %d not found", id)
		return Book{}, fmt.Errorf("book with ID %d not found", id)
	}

	return book, nil
}

func (store *InMemoryBookStore) UpdateBook(ctx context.Context, id int, book Book) (Book, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	book, found := store.books[id]
	if !found {
		log.Printf("book with ID %d not found", id)
		return Book{}, fmt.Errorf("book with ID %d not found", id)
	}
	book.ID = id
	store.books[id] = book

	return book, nil
}

func (store *InMemoryBookStore) DeleteBook(ctx context.Context, id int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, found := store.books[id]
	if !found {
		log.Printf("book with ID %d not found", id)
		return fmt.Errorf("book with ID %d not found", id)
	}

	delete(store.books, id)
	return nil
}

func (store *InMemoryBookStore) GetAllBooks(ctx context.Context) ([]Book, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	all := make([]Book, 0, len(store.books))
	for _, book := range store.books {
		all = append(all, book)
	}
	if len(all) == 0 {
		log.Printf("no books found")
		return nil, fmt.Errorf("no books found")
	}
	return all, nil
}

func (store *InMemoryBookStore) SearchBooks(ctx context.Context, criteria SearchCriteria) ([]Book, error) {
	var result []Book
	for _, book := range store.books {
		if (criteria.Title == "" || book.Title == criteria.Title) &&
			(criteria.Author == author.Author{} || book.Author == criteria.Author) &&
			(criteria.Genre == "" || book.Genre == criteria.Genre) &&
			(criteria.PublishedAt.IsZero() || book.PublishedAt.Equal(criteria.PublishedAt)) &&
			(criteria.PriceRange == [2]float64{0, 0} || (book.Price >= criteria.PriceRange[0] && book.Price <= criteria.PriceRange[1])) {
			result = append(result, book)
		}
	}

	if len(result) == 0 {
		log.Printf("no books found matching criteria")
		return nil, fmt.Errorf("no books found matching criteria")
	}

	return result, nil
}

type Service interface {
	CreateBook(ctx context.Context, b Book) (Book, error)
	GetBook(ctx context.Context, id int) (Book, error)
	UpdateBook(ctx context.Context, id int, b Book) (Book, error)
	DeleteBook(ctx context.Context, id int) error
	GetAllBooks(ctx context.Context) ([]Book, error)
	SearchBooks(ctx context.Context, criteria SearchCriteria) ([]Book, error)

	DecrementStock(ctx context.Context, bookID int, qty int) error
}

type service struct {
	store BookStore
}

func NewService(bookStore BookStore) Service {
	return &service{
		store: bookStore,
	}
}
func (s *service) CreateBook(ctx context.Context, b Book) (Book, error) {
	return s.store.CreateBook(ctx, b)
}
func (s *service) GetBook(ctx context.Context, id int) (Book, error) {
	return s.store.GetBook(ctx, id)
}

func (s *service) UpdateBook(ctx context.Context, id int, b Book) (Book, error) {
	return s.store.UpdateBook(ctx, id, b)
}

func (s *service) DeleteBook(ctx context.Context, id int) error {
	return s.store.DeleteBook(ctx, id)
}

func (s *service) GetAllBooks(ctx context.Context) ([]Book, error) {
	return s.store.GetAllBooks(ctx)
}

func (s *service) SearchBooks(ctx context.Context, criteria SearchCriteria) ([]Book, error) {
	return s.store.SearchBooks(ctx, criteria)
}

func (s *service) DecrementStock(ctx context.Context, bookID int, qty int) error {
	b, err := s.store.GetBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book with ID %d: %w", bookID, err)
	}
	if b.Stock < qty {
		return fmt.Errorf("insufficient stock for book with ID %d; available: %d, requested: %d", bookID, b.Stock, qty)
	}
	b.Stock -= qty

	_, err = s.store.UpdateBook(ctx, bookID, b)
	if err != nil {
		return fmt.Errorf("failed to update book ID %d after decrementing stock: %w", bookID, err)
	}

	return nil
}
func (s *service) AddBookStock(ctx context.Context, bookID int, qty int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		b, err := s.store.GetBook(ctx, bookID)
		if err != nil {
			return fmt.Errorf("failed to get book with ID %d: %w", bookID, err)
		}

		b.Stock += qty
		_, err = s.store.UpdateBook(ctx, bookID, b)
		if err != nil {
			return fmt.Errorf("failed to update book ID %d after incrementing stock: %w", bookID, err)
		}

		return nil
	}
}
