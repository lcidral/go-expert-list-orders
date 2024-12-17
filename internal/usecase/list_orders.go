package usecase

import (
	"go-expert-list-orders/internal/domain"
)

type ListOrdersUseCase struct {
	repo domain.OrderRepository
}

func NewListOrdersUseCase(repo domain.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{repo: repo}
}

func (uc *ListOrdersUseCase) Execute() ([]domain.Order, error) {
	return uc.repo.List()
}
