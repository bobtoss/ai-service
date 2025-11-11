package postgres

import (
	"context"
	"fmt"

	"ai-service/internal/service/document"
)

type DocumentRepository struct {
	db *DB
}

func NewDocumentRepository(db *DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, doc *document.Document) error {
	query := `
		INSERT INTO document (user_id, document_id, documnt_name)
		VALUES ($1, $2, $3)
		RETURNING document_id`
	return r.db.Pool.QueryRow(ctx, query, doc.UserID, doc.DocumentID, doc.DocumentName).
		Scan(&doc.DocumentID)
}

func (r *DocumentRepository) GetByID(ctx context.Context, id string) (*document.Document, error) {
	query := `SELECT user_id, document_id, documnt_name FROM document WHERE document_id = $1`
	row := r.db.Pool.QueryRow(ctx, query, id)
	d := new(document.Document)
	err := row.Scan(&d.UserID, &d.DocumentID, &d.DocumentName)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	return d, nil
}

func (r *DocumentRepository) ListByUser(ctx context.Context, userID string) ([]*document.Document, error) {
	query := `SELECT user_id, document_id, documnt_name FROM document WHERE user_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list document: %w", err)
	}
	defer rows.Close()

	var docs []*document.Document
	for rows.Next() {
		d := new(document.Document)
		if err := rows.Scan(&d.UserID, &d.DocumentID, &d.DocumentName); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM document WHERE document_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
