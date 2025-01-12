package order

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"um6p.ma/final_project/internal/book"
	"um6p.ma/final_project/internal/customer"
	"um6p.ma/final_project/pkg/error"
)

var (
	customerStore customer.CustomerStore = customer.NewCustomerStore()
	bookStore     book.BookStore         = book.NewStore()
	ordStore      OrderStore             = NewOrderStore()
)

func NewOrderStore() *InMemoryOrderStore {
	return &InMemoryOrderStore{
		orders: make(map[int]Order),
		nextID: 1,
	}
}

var orderService = NewService(
	ordStore,
	customerStore,
	bookStore,
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		orders, err := orderService.ListOrders(ctx)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)

	case http.MethodPost:
		var o Order
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			error.WriteJSONError(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		createdOrder, err := orderService.CreateOrder(ctx, o)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdOrder)

	default:
		error.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.URL.Path[len("/orders/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		error.WriteJSONError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ord, err := orderService.GetOrderByID(ctx, id)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ord)

	case http.MethodPut:
		var updated Order
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			error.WriteJSONError(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}
		updated.ID = id

		if err := orderService.UpdateOrder(ctx, id, updated); err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		if err := orderService.DeleteOrder(ctx, id); err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		error.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
