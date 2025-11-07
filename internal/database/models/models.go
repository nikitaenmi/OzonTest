package models

import (
	"github.com/nikitaenmi/OzonTest/internal/domain"
	"time"
)

type Payment struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	ProviderName string  `gorm:"not null"`
	Amount       float64 `gorm:"type:decimal(15,2);not null"`
	Currency     string  `gorm:"size:3;not null"`
	AmountRub    float64 `gorm:"type:decimal(15,2);not null"`
	PaymentDate  time.Time
	CreatedAt    time.Time
}

func (p Payment) FromDomain(domainPayment *domain.Payment) Payment {
	return Payment{
		ProviderName: domainPayment.Provider,
		Amount:       domainPayment.Amount,
		Currency:     domainPayment.Currency,
		AmountRub:    domainPayment.AmountRub,
		PaymentDate:  domainPayment.Date,
	}
}

func (p Payment) ToDomain() domain.Payment {
	return domain.Payment{
		ID:       int(p.ID),
		Provider: p.ProviderName,
		Amount:   p.AmountRub,
		Currency: "RUB",
		Date:     p.PaymentDate,
	}
}
