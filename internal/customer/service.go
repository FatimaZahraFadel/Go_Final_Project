package customer

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type InMemoryCustomerStore struct {
	mu        sync.RWMutex
	customers map[int]Customer
	nextID    int
}

func (store *InMemoryCustomerStore) CreateCustomer(ctx context.Context, customer *Customer) (Customer, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return Customer{}, ctx.Err()
	default:
	}

	for _, existingCustomer := range store.customers {
		if existingCustomer.Name == customer.Name {
			log.Printf("customer with name %s already exists", customer.Name)
			return Customer{}, fmt.Errorf("customer with name %s already exists", customer.Name)
		}
	}

	customer.ID = store.nextID
	store.customers[customer.ID] = *customer
	store.nextID++

	return *customer, nil
}

func (store *InMemoryCustomerStore) GetCustomerByID(ctx context.Context, id int) (Customer, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return Customer{}, ctx.Err()
	default:
	}

	customer, found := store.customers[id]
	if !found {
		log.Printf("customer with ID %d not found", id)
		return Customer{}, fmt.Errorf("customer with ID %d not found", id)
	}

	return customer, nil
}

func (store *InMemoryCustomerStore) UpdateCustomer(ctx context.Context, id int, customer *Customer) (Customer, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	existing, found := store.customers[id]
	if !found {
		return Customer{}, fmt.Errorf("customer with ID %d not found", id)
	}
	existing.Name = customer.Name
	existing.Email = customer.Email
	existing.Address = customer.Address

	store.customers[id] = existing
	return existing, nil
}

func (store *InMemoryCustomerStore) DeleteCustomer(ctx context.Context, id int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, found := store.customers[id]
	if !found {
		log.Printf("customer with ID %d not found", id)
		return fmt.Errorf("customer with ID %d not found", id)
	}

	delete(store.customers, id)
	return nil
}

func (store *InMemoryCustomerStore) GetAllCustomers(ctx context.Context) ([]Customer, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	all := make([]Customer, 0, len(store.customers))
	for _, customer := range store.customers {
		all = append(all, customer)
	}

	if len(all) == 0 {
		log.Printf("no customers found")
		return nil, fmt.Errorf("no customers found")
	}

	return all, nil
}
