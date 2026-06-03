package repository

import (
	"context"
	"warehouse/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) error
	GetStatus(ctx context.Context, orderID string) (domain.Order, error)
	UpdateStatus(ctx context.Context, orderID, status string) error
}
