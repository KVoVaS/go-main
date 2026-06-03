package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	devicesv1 "devices/api/gen/devices/v1"
	"devices/internal/service"
	"devices/internal/transport/grpcsrv"
	"devices/internal/middleware"
)

func main() {
	svc := service.NewDeviceService()
	server := grpcsrv.NewDeviceServer(svc)

	// gRPC сервер с интерсепторами
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthUnaryInterceptor()),
		grpc.StreamInterceptor(middleware.AuthStreamInterceptor()),
	)
	devicesv1.RegisterDeviceServiceServer(s, server)
	reflection.Register(s)

	go func() {
		log.Println("gRPC on :9090")
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// REST шлюз с middleware
	ctx := context.Background()
	mux := runtime.NewServeMux(
		// опция для передачи метаданных от HTTP в gRPC
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs("x-api-key", req.Header.Get("X-API-Key"))
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = devicesv1.RegisterDeviceServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		log.Fatal(err)
	}

	wrappedMux := middleware.AuthHTTPMiddleware(mux)
	log.Println("HTTP on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatal(err)
	}
}
