package postgres

import (
	"context"
	"database/sql"

	"documents/internal/domain"
	postgresgen "documents/internal/repository/postgres/gen" // сгенерированный sqlc-пакет

	"github.com/google/uuid"
)

type DocumentRepo struct {
	db *sql.DB
	q  *postgresgen.Queries
}

func NewDocumentRepo(db *sql.DB) *DocumentRepo {
	return &DocumentRepo{db: db, q: postgresgen.New(db)}
}

func (r *DocumentRepo) CreateUser(ctx context.Context, name string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, name)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{ID: row.ID.String(), Name: row.Name}, nil
}

func (r *DocumentRepo) AddDocument(ctx context.Context, d domain.Document) (domain.Document, error) {
	userID, err := uuid.Parse(d.UserID)
	if err != nil {
		return domain.Document{}, err
	}
	row, err := r.q.AddDocument(ctx, postgresgen.AddDocumentParams{
		Title:  d.Title,
		UserID: userID,
	})
	if err != nil {
		return domain.Document{}, err
	}
	return domain.Document{
		ID:     row.ID.String(),
		Title:  row.Title,
		UserID: row.UserID.String(),
	}, nil
}

func (r *DocumentRepo) ListUserDocuments(ctx context.Context, userID string) ([]domain.Document, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListDocumentsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	docs := make([]domain.Document, 0, len(rows))
	for _, row := range rows {
		docs = append(docs, domain.Document{
			ID:     row.ID.String(),
			Title:  row.Title,
			UserID: row.UserID.String(),
		})
	}
	return docs, nil
}

func (r *DocumentRepo) UserExists(ctx context.Context, userID string) (bool, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}
	return r.q.UserExists(ctx, uid)
}
