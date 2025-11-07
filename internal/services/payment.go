package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikitaenmi/OzonTest/internal/domain"
)

type PaymentService struct {
	repo         domain.PaymentRepo
	conventor    CurrencyConverter
	maxAmountRUB float64
}

func NewPaymentService(repo domain.PaymentRepo, conventor CurrencyConverter, maxAmountRUB float64) *PaymentService {
	return &PaymentService{
		repo:         repo,
		conventor:    conventor,
		maxAmountRUB: maxAmountRUB,
	}
}

func (s *PaymentService) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if payment.Amount < 0 {
		return nil, fmt.Errorf("amount cannot be negative: %.2f", payment.Amount)
	}
	amountRub, err := s.conventor.ToRUB(payment.Amount, payment.Currency, payment.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %v", err)
	}
	payment.AmountRub = amountRub

	if payment.AmountRub > s.maxAmountRUB {
		return nil, errors.New("amount exceeds limit of 15000 RUB")
	}

	createdPayment, err := s.repo.Create(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %v", err)
	}

	return &createdPayment, nil
}
