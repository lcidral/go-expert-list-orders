package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"

	"go-expert-list-orders/internal/domain"
	"go-expert-list-orders/internal/usecase"
)

type OrderHandler struct {
	listOrdersUseCase  *usecase.ListOrdersUseCase
	createOrderUseCase *usecase.CreateOrderUseCase
}

func NewOrderHandler(listUC *usecase.ListOrdersUseCase, createUC *usecase.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{
		listOrdersUseCase:  listUC,
		createOrderUseCase: createUC,
	}
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listOrders(w, r)
	case http.MethodPost:
		h.createOrder(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) listOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.listOrdersUseCase.Execute()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CustomerID  string  `json:"customer_id"`
		Status      string  `json:"status"`
		TotalAmount float64 `json:"total_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	order := &domain.Order{
		ID:          uuid.New().String(), // Gerando novo UUID
		CustomerID:  input.CustomerID,
		Status:      input.Status,
		TotalAmount: input.TotalAmount,
	}

	err := h.createOrderUseCase.Execute(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
