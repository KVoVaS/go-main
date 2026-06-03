package repository

import (
	"context"
	"garage/internal/domain"
)

type CarRepository interface {
	Create(ctx context.Context, car domain.Car) error
	Get(ctx context.Context, vin string) (domain.Car, error)
}
