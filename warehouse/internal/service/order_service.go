package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"warehouse/internal/broker"
	"warehouse/internal/domain"
	"warehouse/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	repo   repository.OrderRepository
	broker broker.MessageBroker
}

func NewOrderService(repo repository.OrderRepository, broker broker.MessageBroker) *OrderService {
	return &OrderService{repo: repo, broker: broker}
}

func (s *OrderService) CreateOrder(ctx context.Context, productID string, quantity int32) (domain.Order, error) {
	order := domain.Order{
		ID:        uuid.New().String(),
		ProductID: productID,
		Quantity:  quantity,
		Status:    "PENDING",
	}
	if err := s.repo.Create(ctx, order); err != nil {
		return domain.Order{}, err
	}

	// Отправляем ID заказа в очередь "orders"
	msg, _ := json.Marshal(map[string]string{"order_id": order.ID})
	if err := s.broker.Publish("orders", msg); err != nil {
		log.Printf("failed to publish message: %v", err)
		// Можно обновить статус на FAILED, если критично
	}
	return order, nil
}

func (s *OrderService) GetOrderStatus(ctx context.Context, orderID string) (domain.Order, error) {
	return s.repo.GetStatus(ctx, orderID)
}

// StartWorker запускает фоновую обработку очереди
func (s *OrderService) StartWorker() {
	s.broker.Subscribe("orders", func(msg []byte) error {
		var data map[string]string
		json.Unmarshal(msg, &data)
		orderID := data["order_id"]

		// Имитация обработки (проверка стока, резервирование)
		time.Sleep(10 * time.Second) // эмуляция работы

		// В реальности: проверить product.stock, уменьшить, обновить статус
		// Для упрощения просто переводим в RESERVED
		err := s.repo.UpdateStatus(context.Background(), orderID, "RESERVED")
		if err != nil {
			log.Printf("failed to update order %s: %v", orderID, err)
			_ = s.repo.UpdateStatus(context.Background(), orderID, "FAILED")
		}
		return err
	})
}
