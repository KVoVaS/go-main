package auth

import (
	"context"
	"errors"
)

var ErrUnauthorized = errors.New("invalid API key")

func ValidateAPIKey(ctx context.Context, key string) error {
	// В реальном проекте — проверка по БД/кэшу
	if key == "secret-api-key" {
		return nil
	}
	return ErrUnauthorized
}
