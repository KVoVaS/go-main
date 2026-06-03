package postgres

import (
	"context"
	"database/sql"
	"warehouse/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, order domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO orders (id, product_id, quantity, status) VALUES ($1, $2, $3, $4)",
		order.ID, order.ProductID, order.Quantity, order.Status)
	return err
}

func (r *OrderRepo) GetStatus(ctx context.Context, orderID string) (domain.Order, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, product_id, quantity, status FROM orders WHERE id = $1", orderID)
	var o domain.Order
	err := row.Scan(&o.ID, &o.ProductID, &o.Quantity, &o.Status)
	return o, err
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, orderID, status string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE orders SET status = $2 WHERE id = $1", orderID, status)
	return err
}
