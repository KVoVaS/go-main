package repository

import (
	"context"
	"documents/internal/domain"
)

type DocumentRepository interface {
	CreateUser(ctx context.Context, name string) (domain.User, error)
	AddDocument(ctx context.Context, doc domain.Document) (domain.Document, error)
	ListUserDocuments(ctx context.Context, userID string) ([]domain.Document, error)
	UserExists(ctx context.Context, userID string) (bool, error)
}
