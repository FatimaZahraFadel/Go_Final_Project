package sales

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"um6p.ma/final_project/internal/order"

	pkgError "um6p.ma/final_project/pkg/error"
)

type InMemorySalesStore struct {
	orders map[int]order.Order
	nextID int
}

func NewSalesStore() *InMemorySalesStore {
	return &InMemorySalesStore{
		orders: make(map[int]order.Order),
		nextID: 1,
	}
}

func (s *InMemorySalesStore) generateSalesReport(ctx context.Context, start, end time.Time) (interface{}, error) {
	placeholder := map[string]interface{}{
		"message": "generateSalesReport placeholder",
		"start":   start.Format(time.RFC3339),
		"end":     end.Format(time.RFC3339),
	}
	return placeholder, nil
}
func SalesReportHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		startStr := q.Get("start")
		endStr := q.Get("end")

		var start, end time.Time
		var err error

		salesService := NewSalesStore()

		if startStr == "" || endStr == "" {
			end = time.Now()
			start = end.Add(-24 * time.Hour)
		} else {
			start, err = time.Parse(time.RFC3339, startStr)
			if err != nil {
				pkgError.WriteJSONError(w, "invalid 'start' time (use RFC3339)", http.StatusBadRequest)
				return
			}
			end, err = time.Parse(time.RFC3339, endStr)
			if err != nil {
				pkgError.WriteJSONError(w, "invalid 'end' time (use RFC3339)", http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		report, err := salesService.generateSalesReport(ctx, start, end)
		if err != nil {
			pkgError.WriteJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)

	default:
		pkgError.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
