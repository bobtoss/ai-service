package repository

import (
	milvusRepo "ai-service/internal/repository/milvus"
	"ai-service/internal/util/config"
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type VectorDB interface {
	GetTopK(k int, search []float32) ([]client.SearchResult, error)
	DeleteDoc(id string) error
	SaveDoc(embeddings [][]float32) error
}

func New(ctx context.Context, cfg *config.Config) (VectorDB, error) {
	return milvusRepo.NewMilvusRepository(ctx, cfg)
}
