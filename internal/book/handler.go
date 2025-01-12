package book

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"um6p.ma/final_project/pkg/error"
)

func NewStore() *InMemoryBookStore {
	return &InMemoryBookStore{
		books:  make(map[int]Book),
		nextID: 1,
	}
}

var store = NewStore()

var svc = NewService(store)

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodPost:
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			error.WriteJSONError(w, "bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		createdBook, err := svc.CreateBook(ctx, book)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBook)

	case http.MethodGet:
		books, err := svc.GetAllBooks(ctx)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)

	default:
		error.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := strconv.Atoi(r.URL.Path[len("/books/"):])
	if err != nil {
		error.WriteJSONError(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		book, err := svc.GetBook(ctx, id)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)

	case http.MethodPut:
		var updatedData Book
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			error.WriteJSONError(w, "bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		updatedData.ID = id
		updatedBook, err := svc.UpdateBook(ctx, id, updatedData)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedBook)

	case http.MethodDelete:
		if err := svc.DeleteBook(ctx, id); err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		error.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
