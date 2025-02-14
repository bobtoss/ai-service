package milvus

import (
	"ai-service/internal/repository/vector"
	"ai-service/internal/util/config"
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Repository struct {
	milvus client.Client
}

func NewMilvusRepository(ctx context.Context, cfg *config.Config) (vector.VectorDB, error) {
	milvus, err := client.NewClient(ctx, client.Config{Address: cfg.Milvus.Host})
	if err != nil {
		return nil, err
	}
	defer milvus.Close()

	return &Repository{milvus: milvus}, nil
}

func (r Repository) GetTopK(k int, search []float32) ([]client.SearchResult, error) {
	sp, _ := entity.NewIndexIvfFlatSearchParam( // NewIndex*SearchParam func
		10, // searchParam
	)
	opt := client.SearchQueryOptionFunc(func(option *client.SearchQueryOption) {
		option.Limit = 3
		option.Offset = 0
		option.ConsistencyLevel = entity.ClStrong
		option.IgnoreGrowing = false
	})
	searchResult, err := milvusCLient.Search(
		context.Background(), // ctx
		"yandex_gpt",         // CollectionName
		[]string{},           // partitionNames
		"",                   // expr
		[]string{"text"},     // outputFields
		[]entity.Vector{entity.FloatVector(search)}, // vectors
		"embedding", // vectorField
		entity.L2,   // metricType
		k,           // topK
		sp,          // sp
		opt,
	)
	for _, sr := range searchResult {
		fmt.Println(sr.IDs)
		fmt.Println(sr.Fields.GetColumn("text"))
	}
	return searchResult, err
}

func (r Repository) DeleteDoc(id string) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) SaveDoc(embeddings [][]float32) error {
	//TODO implement me
	panic("implement me")
}
