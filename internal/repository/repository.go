package repository

import (
	"ai-service/internal/repository/vector"
	milvusRepo "ai-service/internal/repository/vector/milvus"
	"ai-service/internal/util/config"
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Repository struct {
	Vector vector.VectorDB
}

func New(ctx context.Context, cfg *config.Config) (*Repository, client.Client, error) {
	vector, milvus, err := milvusRepo.NewMilvusRepository(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	return &Repository{Vector: vector}, milvus, nil
}
