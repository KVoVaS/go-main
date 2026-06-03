package memory

import (
	"context"
	"errors"
	"sync"

	"garage/internal/domain"
	"garage/internal/repository"
)

type InMemoryCarRepo struct {
	mu   sync.RWMutex
	data map[string]domain.Car
}

// Конструктор, который мы вызываем в main.go
func NewInMemoryCarRepo() repository.CarRepository {
	return &InMemoryCarRepo{
		data: make(map[string]domain.Car),
	}
}

func (r *InMemoryCarRepo) Create(ctx context.Context, car domain.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[car.VIN]; exists {
		return errors.New("car already exists")
	}
	r.data[car.VIN] = car
	return nil
}

func (r *InMemoryCarRepo) Get(ctx context.Context, vin string) (domain.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	car, ok := r.data[vin]
	if !ok {
		return domain.Car{}, errors.New("not found")
	}
	return car, nil
}
