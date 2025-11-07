package server

import (
	"github.com/nikitaenmi/OzonTest/gen"
	"github.com/nikitaenmi/OzonTest/internal/domain"
	"time"
)


func RequestToDomain(req *gen.CreatePaymentRequest) *domain.Payment {
	paymentDate, err := time.Parse("02/01/2006", req.GetDate())
	if err != nil {
		paymentDate = time.Now()
	}

	return &domain.Payment{
		Provider: req.GetProvider(),
		Amount:   req.GetAmount(),
		Date:     paymentDate,
		Currency: req.GetCurrency(),
	}
}

func DomainToResponse(payment *domain.Payment) *gen.Payment {
	dateStr := payment.Date.Format("02/01/2006")

	return &gen.Payment{
		Id:       int32(payment.ID),
		Provider: payment.Provider,
		Amount:   payment.Amount,
		Date:     dateStr,
		Currency: payment.Currency,
	}
}
