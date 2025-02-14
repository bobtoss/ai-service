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
	searchResult, err := r.milvus.Search(
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

func (r Repository) SaveDoc(chunks []string, embeddings [][]float32) error {
	ids := make([]int64, len(embeddings))
	schema := &entity.Schema{
		CollectionName: "yandex_gpt",
		Description:    "testing",
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:       "text",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: false,
				AutoID:     false,
				TypeParams: map[string]string{
					"max_length": "5000",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "256",
				},
			},
		},
		EnableDynamicField: true,
	}
	err := r.milvus.CreateCollection(
		context.Background(), // ctx
		schema,
		2, // shardNum
	)
	if err != nil {
		return err
	}
	idColumn := entity.NewColumnInt64("id", ids)
	textColumn := entity.NewColumnVarChar("text", chunks)
	embeddingColumn := entity.NewColumnFloatVector("embedding", 256, embeddings)
	_, err = r.milvus.Insert(
		context.Background(), // ctx
		"yandex_gpt",         // CollectionName
		"",                   // partitionName
		idColumn,             // columnarData
		textColumn,
		embeddingColumn, // columnarData
	)
	fmt.Println(len(chunks), len(embeddings))
	if err != nil {
		return err
	}
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024,      // ConstructParams
	)
	if err != nil {
		return err
	}
	err = r.milvus.CreateIndex(
		context.Background(), // ctx
		"yandex_gpt",         // CollectionName
		"embedding",          // fieldName
		idx,                  // entity.Index
		false,                // async
	)
	if err != nil {
		return err
	}
	err = r.milvus.LoadCollection(
		context.Background(), // ctx
		"yandex_gpt",         // CollectionName
		false,                // async
	)
	return nil
}
