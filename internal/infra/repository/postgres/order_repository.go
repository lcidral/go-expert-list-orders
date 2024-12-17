package postgres

import (
	"database/sql"
	"go-expert-list-orders/internal/domain"
	"time"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) List() ([]domain.Order, error) {
	orders := []domain.Order{}
	rows, err := r.db.Query(`
        SELECT id, customer_id, status, total_amount, created_at, updated_at 
        FROM orders
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *PostgresOrderRepository) Create(order *domain.Order) error {
	stmt, err := r.db.Prepare(`
        INSERT INTO orders (id, customer_id, status, total_amount, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	_, err = stmt.Exec(
		order.ID,
		order.CustomerID,
		order.Status,
		order.TotalAmount,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}
