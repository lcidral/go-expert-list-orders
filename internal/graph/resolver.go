package graph

import (
	"go-expert-list-orders/internal/usecase"
)

type Resolver struct {
	listOrdersUseCase *usecase.ListOrdersUseCase
}

func NewResolver(listOrdersUseCase *usecase.ListOrdersUseCase) *Resolver {
	return &Resolver{
		listOrdersUseCase: listOrdersUseCase,
	}
}
