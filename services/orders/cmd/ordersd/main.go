package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Product struct {
	ID            string `json:"id"`
	SKU           string `json:"sku"`
	Title         string `json:"title"`
	ProductType   string `json:"product_type"`
	Niche         string `json:"niche"`
	PriceCents    int64  `json:"price_cents"`
	DownloadAsset string `json:"download_asset"`
}

type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	SKU         string    `json:"sku"`
	Customer    Customer  `json:"customer"`
	GrossCents  int64     `json:"gross_cents"`
	FeeCents    int64     `json:"fee_cents"`
	NetCents    int64     `json:"net_cents"`
	Status      string    `json:"status"`
	DeliveredAt time.Time `json:"delivered_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store struct {
	mu       sync.RWMutex
	products map[string]Product
	orders   map[string]Order
	nextID   int64
}

func NewStore() *Store {
	product := Product{
		ID:            "prod_dark_botanical_wallpaper_pack",
		SKU:           "DB-WALLPAPER-001",
		Title:         "Dark Botanical Wallpaper Pack",
		ProductType:   "wallpaper",
		Niche:         "dark-botanical",
		PriceCents:    900,
		DownloadAsset: "dark-botanical-wallpaper-pack.zip",
	}

	return &Store{
		products: map[string]Product{product.ID: product},
		orders:   make(map[string]Order),
		nextID:   1,
	}
}

func (s *Store) CreateOrder(productID string, customer Customer, now time.Time) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[productID]
	if !ok {
		return Order{}, errors.New("product not found")
	}

	if strings.TrimSpace(customer.Name) == "" {
		return Order{}, errors.New("customer name is required")
	}

	feeCents := calculateEtsyLikeFee(product.PriceCents)
	order := Order{
		ID:          "ord_" + now.UTC().Format("20060102150405") + "_" + formatOrderSequence(s.nextID),
		ProductID:   product.ID,
		SKU:         product.SKU,
		Customer:    customer,
		GrossCents:  product.PriceCents,
		FeeCents:    feeCents,
		NetCents:    product.PriceCents - feeCents,
		Status:      "delivered",
		DeliveredAt: now.UTC(),
		CreatedAt:   now.UTC(),
	}

	s.nextID++
	s.orders[order.ID] = order
	return order, nil
}

func (s *Store) ListOrders() []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

func calculateEtsyLikeFee(priceCents int64) int64 {
	listingFeeCents := int64(20)
	transactionFeeCents := priceCents * 65 / 1000
	paymentFeeCents := priceCents*3/100 + 25
	return listingFeeCents + transactionFeeCents + paymentFeeCents
}

func formatOrderSequence(sequence int64) string {
	return fmt.Sprintf("%04d", sequence)
}

type Server struct {
	store *Store
	log   *slog.Logger
}

func (s Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /orders", s.listOrders)
	mux.HandleFunc("POST /orders/sample-john-doe", s.createJohnDoeSampleSale)
	return requestLogger(s.log, mux)
}

func (s Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s Server) listOrders(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"orders": s.store.ListOrders()})
}

func (s Server) createJohnDoeSampleSale(w http.ResponseWriter, r *http.Request) {
	order, err := s.store.CreateOrder("prod_dark_botanical_wallpaper_pack", Customer{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}, time.Now())
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func requestLogger(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("failed to write json response", "error", err)
	}
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := Server{
		store: NewStore(),
		log:   log,
	}

	addr := ":8080"
	log.Info("orders service listening", "addr", addr)
	if err := http.ListenAndServe(addr, server.routes()); err != nil {
		log.Error("orders service failed", "error", err)
		os.Exit(1)
	}
}
