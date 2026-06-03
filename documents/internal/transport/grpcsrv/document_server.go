package grpcsrv

import (
	"context"
	documentsv1 "documents/api/gen/documents/v1"
	"documents/internal/domain"
	"documents/internal/service"
)

type DocumentServer struct {
	documentsv1.UnimplementedDocumentServiceServer
	svc *service.DocumentService
}

func NewDocumentServer(svc *service.DocumentService) *DocumentServer {
	return &DocumentServer{svc: svc}
}
func (s *DocumentServer) CreateUser(ctx context.Context, req *documentsv1.CreateUserRequest) (*documentsv1.CreateUserResponse, error) {
	user, err := s.svc.CreateUser(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &documentsv1.CreateUserResponse{User: &documentsv1.User{Id: user.ID, Name: user.Name}}, nil
}
func (s *DocumentServer) AddDocumentToUser(ctx context.Context, req *documentsv1.AddDocumentRequest) (*documentsv1.AddDocumentResponse, error) {
	doc, err := s.svc.AddDocument(ctx, req.UserId, domain.Document{Title: req.Document.Title})
	if err != nil {
		return nil, err
	}
	return &documentsv1.AddDocumentResponse{Document: &documentsv1.Document{Id: doc.ID, Title: doc.Title, UserId: doc.UserID}}, nil
}
func (s *DocumentServer) ListUserDocuments(ctx context.Context, req *documentsv1.ListUserDocumentsRequest) (*documentsv1.ListUserDocumentsResponse, error) {
	docs, err := s.svc.ListUserDocuments(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var pbDocs []*documentsv1.Document
	for _, d := range docs {
		pbDocs = append(pbDocs, &documentsv1.Document{Id: d.ID, Title: d.Title, UserId: d.UserID})
	}
	return &documentsv1.ListUserDocumentsResponse{Documents: pbDocs}, nil
}
