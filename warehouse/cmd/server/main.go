package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	warehousev1 "warehouse/api/gen/warehouse/v1"
	"warehouse/internal/broker/memory"
	"warehouse/internal/repository/postgres"
	"warehouse/internal/service"
	"warehouse/internal/transport/grpcsrv"
)

func main() {
	// Подключение к БД
	connStr := "postgres://postgres:secret@localhost:5432/warehouse?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация компонентов
	repo := postgres.NewOrderRepo(db)
	broker := memory.NewMemoryBroker()
	svc := service.NewOrderService(repo, broker)

	// Запускаем воркер
	svc.StartWorker()

	// gRPC сервер
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	warehousev1.RegisterWarehouseServiceServer(s, grpcsrv.NewWarehouseServer(svc))
	reflection.Register(s)

	go func() {
		log.Println("gRPC on :9090")
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// REST шлюз
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = warehousev1.RegisterWarehouseServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("HTTP on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
