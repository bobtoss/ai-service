package repository

import (
	"ai-service/internal/repository/vector"
	milvusRepo "ai-service/internal/repository/vector/milvus"
	"ai-service/internal/util/config"
	"context"
)

type Repository struct {
	Vector vector.VectorDB
}

func New(ctx context.Context, cfg *config.Config) (*Repository, error) {
	vector, err := milvusRepo.NewMilvusRepository(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{Vector: vector}, nil
}
