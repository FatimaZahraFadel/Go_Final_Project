package sales

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"um6p.ma/final_project/internal/order"
)

type Service interface {
	StartPeriodicReportGeneration(ctx context.Context)
	Stop()
}

type service struct {
	orderStore order.OrderStore
	salesStore SalesStore

	ticker  *time.Ticker
	stopCh  chan struct{}
	wg      sync.WaitGroup
	running bool
}

func NewService(oStore order.OrderStore, sStore SalesStore) Service {
	return &service{
		orderStore: oStore,
		salesStore: sStore,
		stopCh:     make(chan struct{}),
	}
}

func (s *service) StartPeriodicReportGeneration(ctx context.Context) {
	if s.running {
		return
	}
	s.running = true

	s.ticker = time.NewTicker(24 * time.Hour)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ticker.C:
				if _, err := s.generateSalesReport(ctx, time.Now().Add(-24*time.Hour), time.Now()); err != nil {
					log.Printf("SalesService Failed to generate report: %v\n", err)
				} else {
					log.Println("SalesService Sales report generated successfully.")
				}

			case <-s.stopCh:
				s.ticker.Stop()
				return

			case <-ctx.Done():
				log.Println("SalesService context cancelled")
				s.Stop()
				return
			}
		}
	}()
}

func (s *service) Stop() {
	if !s.running {
		return
	}
	close(s.stopCh)
	s.wg.Wait()
	s.running = false
}

func (s *service) generateSalesReport(ctx context.Context, start, end time.Time) (SalesReport, error) {
	orders, err := s.orderStore.GetOrdersInTimeRange(ctx, start, end)
	if err != nil {
		return SalesReport{}, fmt.Errorf("getOrdersInTimeRange: %w", err)
	}

	var totalRevenue float64
	totalOrders := len(orders)
	bookSalesMap := make(map[int]int)

	for _, o := range orders {
		select {
		case <-ctx.Done():
			return SalesReport{}, ctx.Err()
		default:
			totalRevenue += o.TotalPrice
			for _, item := range o.Items {
				bookSalesMap[item.Book.ID] += item.Quantity
			}
		}
	}

	var topSelling []BookSales
	for _, o := range orders {
		for _, item := range o.Items {
			qty := bookSalesMap[item.Book.ID]
			topSelling = append(topSelling, BookSales{
				Book:     item.Book,
				Quantity: qty,
			})
		}
	}

	report := SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSelling,
	}
	if err := s.salesStore.RecordSale(ctx, BookSales{}); err != nil {
		log.Printf("RecordSale stub. Replace with real logic. err=%v", err)
	}

	return report, nil
}
