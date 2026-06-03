package main

import (
	"context"
	"database/sql"
	documentsv1 "documents/api/gen/documents/v1"
	"documents/internal/repository/postgres"
	"documents/internal/service"
	"documents/internal/transport/grpcsrv"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/documents?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	repo := postgres.NewDocumentRepo(db)
	svc := service.NewDocumentService(repo)
	srv := grpcsrv.NewDocumentServer(svc)
	lis, _ := net.Listen("tcp", ":9090")
	s := grpc.NewServer()
	documentsv1.RegisterDocumentServiceServer(s, srv)
	reflection.Register(s)
	go s.Serve(lis)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	documentsv1.RegisterDocumentServiceHandlerFromEndpoint(context.Background(), mux, "localhost:9090", opts)
	http.ListenAndServe(":8080", mux)
}
