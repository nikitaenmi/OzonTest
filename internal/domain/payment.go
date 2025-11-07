package domain

import (
	"context"
	"time"
)

type Payment struct {
	ID        int
	Provider  string
	Amount    float64
	Currency  string
	AmountRub float64
	Date      time.Time
}

type PaymentRepo interface {
	Create(*Payment) (Payment, error)
}

type PaymentService interface {
	Create(context.Context, *Payment) (*Payment, error)
}
