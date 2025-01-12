package customer

import (
	"encoding/json"
	"net/http"
	"strconv"

	"um6p.ma/final_project/pkg/error"
)

func NewCustomerStore() *InMemoryCustomerStore {
	return &InMemoryCustomerStore{
		customers: make(map[int]Customer),
		nextID:    1,
	}
}

var customerStore = NewCustomerStore()

func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		var customer Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			error.WriteJSONError(w, "bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		createdCustomer, err := customerStore.CreateCustomer(ctx, &customer)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdCustomer)

	case http.MethodGet:
		customers, err := customerStore.GetAllCustomers(ctx)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)

	default:
		error.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func CustomerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(r.URL.Path[len("/customers/"):])
	if err != nil {
		error.WriteJSONError(w, "invalid customer ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		customer, err := customerStore.GetCustomerByID(ctx, id)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customer)

	case http.MethodPut:
		var updatedData Customer
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			error.WriteJSONError(w, "bad request: invalid JSON", http.StatusBadRequest)
			return
		}
		updatedData.ID = id
		updatedCustomer, err := customerStore.UpdateCustomer(ctx, id, &updatedData)
		if err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedCustomer)

	case http.MethodDelete:
		if err := customerStore.DeleteCustomer(ctx, id); err != nil {
			error.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		error.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
