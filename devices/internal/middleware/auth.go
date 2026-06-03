package middleware

import (
	"context"
	"devices/internal/auth"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		keys := md["x-api-key"]
		if len(keys) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing API key")
		}
		if err := auth.ValidateAPIKey(ctx, keys[0]); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return handler(ctx, req)
	}
}

// Stream интерсептор (для MonitorReadings)
func AuthStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.Unauthenticated, "missing metadata")
		}
		keys := md["x-api-key"]
		if len(keys) == 0 {
			return status.Error(codes.Unauthenticated, "missing API key")
		}
		if err := auth.ValidateAPIKey(ss.Context(), keys[0]); err != nil {
			return status.Error(codes.PermissionDenied, err.Error())
		}
		return handler(srv, ss)
	}
}

func AuthHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, "missing API key", http.StatusUnauthorized)
			return
		}
		if err := auth.ValidateAPIKey(r.Context(), key); err != nil {
			http.Error(w, "invalid API key", http.StatusForbidden)
			return
		}
		// Пробрасываем API-ключ в gRPC-метаданные, чтобы шлюз передал его дальше
		ctx := metadata.AppendToOutgoingContext(r.Context(), "x-api-key", key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
