// internal/usecase/create_order.go
package usecase

import (
	"go-expert-list-orders/internal/domain"
)

type CreateOrderUseCase struct {
	repo domain.OrderRepository
}

func NewCreateOrderUseCase(repo domain.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		repo: repo,
	}
}

func (uc *CreateOrderUseCase) Execute(order *domain.Order) error {
	return uc.repo.Create(order)
}
