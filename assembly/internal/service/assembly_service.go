package service

import (
	"assembly/internal/domain"
	"assembly/internal/repository"
	"context"
)

type AssemblyService struct {
	repo repository.AssemblyRepository
}

func NewAssemblyService(repo repository.AssemblyRepository) *AssemblyService {
	return &AssemblyService{repo: repo}
}

func (s *AssemblyService) CreateEngine(ctx context.Context, e domain.Engine) (domain.Engine, error) {
	return s.repo.CreateEngine(ctx, e)
}

func (s *AssemblyService) CreateTransmission(ctx context.Context, t domain.Transmission) (domain.Transmission, error) {
	return s.repo.CreateTransmission(ctx, t)
}

func (s *AssemblyService) AssembleCar(ctx context.Context, vin, brand string, year int32, engineID, transID string) (domain.CarSpec, error) {
	car := domain.Car{VIN: vin, Brand: brand, Year: year}
	return s.repo.AssembleCar(ctx, car, engineID, transID)
}

func (s *AssemblyService) GetCarSpec(ctx context.Context, vin string) (domain.CarSpec, error) {
	return s.repo.GetCarSpec(ctx, vin)
}
