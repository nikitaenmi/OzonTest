package handler

import (
	"context"

	payment "github.com/nikitaenmi/OzonTest/gen"
	"github.com/nikitaenmi/OzonTest/internal/domain"
	"github.com/nikitaenmi/OzonTest/internal/server"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment struct {
	payment.UnimplementedPaymentServiceServer
	service domain.PaymentService
}

func NewPaymentHandler(service domain.PaymentService) *Payment {
	return &Payment{service: service}
}

func (h *Payment) Create(ctx context.Context, req *payment.CreatePaymentRequest) (*payment.Payment, error) {
	domainPayment := server.RequestToDomain(req)
	payment, err := h.service.Create(ctx, domainPayment)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	response := server.DomainToResponse(payment)
	return response, nil
}
