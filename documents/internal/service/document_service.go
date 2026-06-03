package service

import (
	"context"
	"documents/internal/domain"
	"documents/internal/repository"
	"errors"
)

type DocumentService struct{ repo repository.DocumentRepository }

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo}
}
func (s *DocumentService) CreateUser(ctx context.Context, name string) (domain.User, error) {
	return s.repo.CreateUser(ctx, name)
}
func (s *DocumentService) AddDocument(ctx context.Context, userID string, doc domain.Document) (domain.Document, error) {
	exists, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return domain.Document{}, err
	}
	if !exists {
		return domain.Document{}, errors.New("user not found")
	}
	return s.repo.AddDocument(ctx, domain.Document{Title: doc.Title, UserID: userID})
}
func (s *DocumentService) ListUserDocuments(ctx context.Context, userID string) ([]domain.Document, error) {
	return s.repo.ListUserDocuments(ctx, userID)
}
