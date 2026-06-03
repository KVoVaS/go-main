package repository

import (
	"assembly/internal/domain"
	"context"
)

type AssemblyRepository interface {
	CreateEngine(ctx context.Context, e domain.Engine) (domain.Engine, error)
	CreateTransmission(ctx context.Context, t domain.Transmission) (domain.Transmission, error)
	AssembleCar(ctx context.Context, car domain.Car, engineID, transID string) (domain.CarSpec, error)
	GetCarSpec(ctx context.Context, vin string) (domain.CarSpec, error)
}
