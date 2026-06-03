package service

import (
	"context"

	"garage/internal/domain"
	"garage/internal/repository"
)

type CarService struct {
	repo repository.CarRepository
}

// Конструктор
func NewCarService(repo repository.CarRepository) *CarService {
	return &CarService{repo: repo}
}

func (s *CarService) CreateCar(ctx context.Context, car domain.Car) (domain.Car, error) {
	if err := s.repo.Create(ctx, car); err != nil {
		return domain.Car{}, err
	}
	return car, nil
}

func (s *CarService) GetCar(ctx context.Context, vin string) (domain.Car, error) {
	return s.repo.Get(ctx, vin)
}
