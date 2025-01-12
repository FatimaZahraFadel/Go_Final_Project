package order

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"um6p.ma/final_project/internal/book"
	"um6p.ma/final_project/internal/customer"
)

type InMemoryOrderStore struct {
	mu     sync.RWMutex
	orders map[int]Order
	nextID int
}

func (store *InMemoryOrderStore) GetByID(ctx context.Context, id int) (Order, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return Order{}, ctx.Err()
	default:
	}

	order, found := store.orders[id]
	if !found {
		log.Printf("order with ID %d not found", id)
		return Order{}, fmt.Errorf("order with ID %d not found", id)
	}

	return order, nil
}

func (store *InMemoryOrderStore) Update(ctx context.Context, id int, order Order) (Order, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return Order{}, ctx.Err()
	default:
	}

	_, found := store.orders[id]
	if !found {
		log.Printf("order with ID %d not found", id)
		return Order{}, fmt.Errorf("order with ID %d not found", id)
	}
	order.ID = id
	store.orders[id] = order

	return order, nil
}

func (store *InMemoryOrderStore) Delete(ctx context.Context, id int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, found := store.orders[id]
	if !found {
		log.Printf("order with ID %d not found", id)
		return fmt.Errorf("order with ID %d not found", id)
	}

	delete(store.orders, id)
	return nil
}

func (store *InMemoryOrderStore) List(ctx context.Context) ([]Order, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	all := make([]Order, 0, len(store.orders))
	for _, order := range store.orders {
		all = append(all, order)
	}
	if len(all) == 0 {
		log.Printf("no orders found")
		return nil, fmt.Errorf("no orders found")
	}
	return all, nil
}

func (store *InMemoryOrderStore) Create(ctx context.Context, order Order) (Order, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	select {
	case <-ctx.Done():
		return Order{}, ctx.Err()
	default:
	}

	order.ID = store.nextID
	store.nextID++
	store.orders[order.ID] = order

	log.Printf("Order with ID %d created", order.ID)
	return order, nil
}

type Service interface {
	CreateOrder(ctx context.Context, o Order) (Order, error)
	GetOrderByID(ctx context.Context, id int) (Order, error)
	UpdateOrder(ctx context.Context, id int, o Order) error
	DeleteOrder(ctx context.Context, id int) error
	ListOrders(ctx context.Context) ([]Order, error)
}

type service struct {
	store         OrderStore
	customerStore customer.CustomerStore
	bookStore     book.BookStore
}

func NewService(orderStore OrderStore, cStore customer.CustomerStore, bStore book.BookStore) Service {
	return &service{
		store:         orderStore,
		customerStore: cStore,
		bookStore:     bStore,
	}
}
func (s *service) CreateOrder(ctx context.Context, o Order) (Order, error) {
	existingCustomer, err := s.customerStore.GetCustomerByID(ctx, o.Customer.ID)
	if err != nil {
		return Order{}, fmt.Errorf("customer with ID %d not found: %w", o.Customer.ID, err)
	}
	o.Customer = existingCustomer

	for i, item := range o.Items {
		b, err := s.bookStore.GetBook(ctx, item.Book.ID)
		if err != nil {
			return Order{}, fmt.Errorf("book with ID %d not found: %w", item.Book.ID, err)
		}
		if b.Stock < item.Quantity {
			return Order{}, fmt.Errorf("insufficient stock for book with ID %d (available=%d, needed=%d)", b.ID, b.Stock, item.Quantity)
		}
		b.Stock -= item.Quantity
		if _, err := s.bookStore.UpdateBook(ctx, b.ID, b); err != nil {
			return Order{}, fmt.Errorf("failed to update book with ID %d: %w", b.ID, err)
		}
		o.Items[i].Book = b
	}
	errChan := make(chan error, len(o.Items))
	var wg sync.WaitGroup
	var totalPrice float64
	var mu sync.Mutex

	for _, item := range o.Items {
		wg.Add(1)
		go func(it OrderItem) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				cost := it.Book.Price * float64(it.Quantity)

				mu.Lock()
				totalPrice += cost
				mu.Unlock()
			}
		}(item)
	}
	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return Order{}, <-errChan
	}

	o.CreatedAt = time.Now()
	o.TotalPrice = totalPrice
	o.Status = "Pending"

	newOrder, err := s.store.Create(ctx, o)
	if err != nil {
		return Order{}, fmt.Errorf("failed to create order: %w", err)
	}
	return newOrder, nil
}

func (s *service) GetOrderByID(ctx context.Context, id int) (Order, error) {
	return s.store.GetByID(ctx, id)
}
func (s *service) UpdateOrder(ctx context.Context, id int, o Order) error {
	_, err := s.store.Update(ctx, id, o)
	return err
}
func (s *service) DeleteOrder(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}
func (s *service) ListOrders(ctx context.Context) ([]Order, error) {
	return s.store.List(ctx)
}
func (store *InMemoryOrderStore) GetOrdersInTimeRange(ctx context.Context, start, end time.Time) ([]Order, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var result []Order
	for _, o := range store.orders {
		if o.CreatedAt.After(start) && o.CreatedAt.Before(end) {
			result = append(result, o)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no orders found between %s and %s", start, end)
	}
	return result, nil
}
