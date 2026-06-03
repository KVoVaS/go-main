package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	garagev1 "garage/api/gen/garage/v1" // сгенерированный код
	"garage/internal/repository/memory"
	"garage/internal/service"
	"garage/internal/transport/grpcsrv"
)

func main() {
	// Инициализация зависимостей
	repo := memory.NewInMemoryCarRepo()
	svc := service.NewCarService(repo)
	carServer := grpcsrv.NewCarServer(svc)

	// gRPC сервер на порту 9090
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	garagev1.RegisterCarServiceServer(s, carServer)
	reflection.Register(s)
	go func() {
		log.Println("gRPC server listening on :9090")
		log.Fatal(s.Serve(lis))
	}()

	// REST шлюз на порту 8080
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = garagev1.RegisterCarServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
