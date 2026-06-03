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

	assemblyv1 "assembly/api/gen/assembly/v1"
	"assembly/internal/repository/postgres"
	"assembly/internal/service"
	"assembly/internal/transport/grpcsrv"
)

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/assembly?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewAssemblyRepo(db)
	svc := service.NewAssemblyService(repo)
	srv := grpcsrv.NewAssemblyServer(svc)

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	assemblyv1.RegisterAssemblyServiceServer(s, srv)
	reflection.Register(s)

	go func() {
		log.Println("gRPC на :9090")
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = assemblyv1.RegisterAssemblyServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("HTTP на :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
