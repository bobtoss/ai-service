package vector

import "github.com/milvus-io/milvus-sdk-go/v2/client"

type VectorDB interface {
	GetTopK(k int, search []float32) ([]client.SearchResult, error)
	DeleteDoc(id string) error
	SaveDoc(embeddings [][]float32) error
}
