package main

import (
	"testing"
	"time"
)

func TestCreateOrderForSampleDigitalProduct(t *testing.T) {
	store := NewStore()
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)

	order, err := store.CreateOrder("prod_dark_botanical_wallpaper_pack", Customer{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}, now)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if order.ID != "ord_20260529120000_0001" {
		t.Fatalf("order ID = %q, want %q", order.ID, "ord_20260529120000_0001")
	}
	if order.Customer.Name != "John Doe" {
		t.Fatalf("customer name = %q, want John Doe", order.Customer.Name)
	}
	if order.GrossCents != 900 {
		t.Fatalf("gross cents = %d, want 900", order.GrossCents)
	}
	if order.NetCents != 770 {
		t.Fatalf("net cents = %d, want 770", order.NetCents)
	}
	if order.Status != "delivered" {
		t.Fatalf("status = %q, want delivered", order.Status)
	}
}

func TestCreateOrderRejectsUnknownProduct(t *testing.T) {
	store := NewStore()

	_, err := store.CreateOrder("missing", Customer{Name: "John Doe"}, time.Now())
	if err == nil {
		t.Fatal("CreateOrder returned nil error for missing product")
	}
}

func TestCalculateEtsyLikeFee(t *testing.T) {
	got := calculateEtsyLikeFee(900)
	want := int64(130)

	if got != want {
		t.Fatalf("calculateEtsyLikeFee(900) = %d, want %d", got, want)
	}
}
