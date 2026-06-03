package service

import (
	"log"
	"sync"

	"devices/internal/domain"
)

type DeviceService struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan domain.Reading]struct{} // deviceID -> set of channels
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		subscribers: make(map[string]map[chan domain.Reading]struct{}),
	}
}

// PublishReading принимает показание, сохраняет (опционально) и рассылает подписчикам
func (s *DeviceService) PublishReading(reading domain.Reading) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if chs, ok := s.subscribers[reading.DeviceID]; ok {
		for ch := range chs {
			select {
			case ch <- reading:
			default:
				log.Printf("dropping reading for slow subscriber on device %s", reading.DeviceID)
			}
		}
	}
}

// Subscribe подписывает новый канал на показания устройства
func (s *DeviceService) Subscribe(deviceID string, ch chan domain.Reading) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subscribers[deviceID] == nil {
		s.subscribers[deviceID] = make(map[chan domain.Reading]struct{})
	}
	s.subscribers[deviceID][ch] = struct{}{}
}

// Unsubscribe удаляет канал из подписчиков
func (s *DeviceService) Unsubscribe(deviceID string, ch chan domain.Reading) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if chs, ok := s.subscribers[deviceID]; ok {
		delete(chs, ch)
		if len(chs) == 0 {
			delete(s.subscribers, deviceID)
		}
	}
}
