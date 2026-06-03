package grpcsrv

import (
	"context"

	garagev1 "garage/api/gen/garage/v1"
	"garage/internal/domain"
	"garage/internal/service"
)

type CarServer struct {
	garagev1.UnimplementedCarServiceServer
	svc *service.CarService
}

// Конструктор
func NewCarServer(svc *service.CarService) *CarServer {
	return &CarServer{svc: svc}
}

func (s *CarServer) CreateCar(ctx context.Context, req *garagev1.CreateCarRequest) (*garagev1.CreateCarResponse, error) {
	car, err := s.svc.CreateCar(ctx, domain.Car{
		VIN:   req.Vin,
		Brand: req.Brand,
		Year:  req.Year,
	})
	if err != nil {
		return nil, err
	}
	return &garagev1.CreateCarResponse{
		Car: &garagev1.Car{
			Vin:   car.VIN,
			Brand: car.Brand,
			Year:  car.Year,
		},
	}, nil
}

func (s *CarServer) GetCar(ctx context.Context, req *garagev1.GetCarRequest) (*garagev1.GetCarResponse, error) {
	car, err := s.svc.GetCar(ctx, req.Vin)
	if err != nil {
		return nil, err
	}
	return &garagev1.GetCarResponse{
		Car: &garagev1.Car{
			Vin:   car.VIN,
			Brand: car.Brand,
			Year:  car.Year,
		},
	}, nil
}
