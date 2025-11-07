package repositories

import (
	"github.com/nikitaenmi/OzonTest/internal/database/models"
	"github.com/nikitaenmi/OzonTest/internal/domain"

	"gorm.io/gorm"
)

type Payment struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Payment {
	return &Payment{db: db}
}

func (r *Payment) Create(payment *domain.Payment) (domain.Payment, error) {
	var model models.Payment
	paymentModel := model.FromDomain(payment)

	result := r.db.Create(&paymentModel)
	if result.Error != nil {
		return domain.Payment{}, result.Error
	}

	return paymentModel.ToDomain(), nil
}
