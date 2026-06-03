package grpcsrv

import (
	"context"
	"time"

	devicesv1 "devices/api/gen/devices/v1"
	"devices/internal/domain"
	"devices/internal/service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DeviceServer struct {
	devicesv1.UnimplementedDeviceServiceServer
	svc *service.DeviceService
}

func NewDeviceServer(svc *service.DeviceService) *DeviceServer {
	return &DeviceServer{svc: svc}
}

func (s *DeviceServer) PublishReading(ctx context.Context, req *devicesv1.PublishReadingRequest) (*devicesv1.PublishReadingResponse, error) {
	reading := domain.Reading{
		DeviceID:  req.DeviceId,
		Value:     req.Value,
		Timestamp: time.Now(),
	}
	s.svc.PublishReading(reading)
	return &devicesv1.PublishReadingResponse{}, nil
}

func (s *DeviceServer) MonitorReadings(req *devicesv1.MonitorRequest, stream devicesv1.DeviceService_MonitorReadingsServer) error {
	ch := make(chan domain.Reading, 10)
	s.svc.Subscribe(req.DeviceId, ch)
	defer s.svc.Unsubscribe(req.DeviceId, ch)

	for reading := range ch {
		ts := timestamppb.New(reading.Timestamp)
		if err := stream.Send(&devicesv1.Reading{
			DeviceId:  reading.DeviceID,
			Value:     reading.Value,
			Timestamp: ts,
		}); err != nil {
			return err
		}
	}
	return nil
}
