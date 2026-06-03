package grpcsrv

import (
	assemblyv1 "assembly/api/gen/assembly/v1"
	"assembly/internal/domain"
	"assembly/internal/service"
	"context"
)

type AssemblyServer struct {
	assemblyv1.UnimplementedAssemblyServiceServer
	svc *service.AssemblyService
}

func NewAssemblyServer(svc *service.AssemblyService) *AssemblyServer {
	return &AssemblyServer{svc: svc}
}

func (s *AssemblyServer) CreateEngine(ctx context.Context, req *assemblyv1.CreateEngineRequest) (*assemblyv1.Engine, error) {
	e, err := s.svc.CreateEngine(ctx, domain.Engine{ID: req.Id, Horsepower: req.Horsepower})
	if err != nil {
		return nil, err
	}
	return &assemblyv1.Engine{Id: e.ID, Horsepower: e.Horsepower}, nil
}

func (s *AssemblyServer) CreateTransmission(ctx context.Context, req *assemblyv1.CreateTransmissionRequest) (*assemblyv1.Transmission, error) {
	t, err := s.svc.CreateTransmission(ctx, domain.Transmission{ID: req.Id, Type: req.Type})
	if err != nil {
		return nil, err
	}
	return &assemblyv1.Transmission{Id: t.ID, Type: t.Type}, nil
}

func (s *AssemblyServer) AssembleCar(ctx context.Context, req *assemblyv1.AssembleCarRequest) (*assemblyv1.AssembleCarResponse, error) {
	spec, err := s.svc.AssembleCar(ctx, req.Vin, req.Brand, req.Year, req.EngineId, req.TransmissionId)
	if err != nil {
		return nil, err
	}
	return &assemblyv1.AssembleCarResponse{
		Spec: &assemblyv1.CarSpec{
			Car:          &assemblyv1.Car{Vin: spec.Car.VIN, Brand: spec.Car.Brand, Year: spec.Car.Year},
			Engine:       &assemblyv1.Engine{Id: spec.Engine.ID, Horsepower: spec.Engine.Horsepower},
			Transmission: &assemblyv1.Transmission{Id: spec.Transmission.ID, Type: spec.Transmission.Type},
		},
	}, nil
}

func (s *AssemblyServer) GetCarSpec(ctx context.Context, req *assemblyv1.GetCarSpecRequest) (*assemblyv1.CarSpec, error) {
	spec, err := s.svc.GetCarSpec(ctx, req.Vin)
	if err != nil {
		return nil, err
	}
	return &assemblyv1.CarSpec{
		Car:          &assemblyv1.Car{Vin: spec.Car.VIN, Brand: spec.Car.Brand, Year: spec.Car.Year},
		Engine:       &assemblyv1.Engine{Id: spec.Engine.ID, Horsepower: spec.Engine.Horsepower},
		Transmission: &assemblyv1.Transmission{Id: spec.Transmission.ID, Type: spec.Transmission.Type},
	}, nil
}
