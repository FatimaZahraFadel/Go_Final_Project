package author

import (
	"encoding/json"
	"net/http"
	"strconv"

	"um6p.ma/final_project/pkg/error"
)

func NewStore() *InMemoryAuthorStore {
	return &InMemoryAuthorStore{
		authors: make(map[int]Author),
		nextID:  1,
	}
}

var store = NewStore()

func AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		var author Author
		if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
			error.WriteJSONError(w, "bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		createdAuthor, err := store.CreateAuthor(ctx, author)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdAuthor)

	case http.MethodGet:
		authors, err := store.ListAuthors(ctx)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authors)

	default:
		error.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func AuthorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.URL.Path[len("/authors/"):])
	if err != nil {
		error.WriteJSONError(w, "invalid author ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		author, err := store.GetAuthorByID(ctx, id)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(author)

	case http.MethodPut:
		var updatedData Author
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			error.WriteJSONError(w, "bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		updatedAuthor := store.UpdateAuthor(ctx, id, updatedData)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedAuthor)

	case http.MethodDelete:
		if err := store.DeleteAuthor(ctx, id); err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		error.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
