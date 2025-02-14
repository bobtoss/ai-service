package vector

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type VectorDB interface {
	GetTopK(ctx context.Context, k int, search []float32) ([]client.SearchResult, error)
	DeleteDoc(ctx context.Context, id string) error
	SaveDoc(ctx context.Context, chunks []string, embeddings [][]float32) error
}
