package milvus

import (
	"ai-service/internal/repository/vector"
	"ai-service/internal/util/config"
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"strings"
)

type Repository struct {
	milvus client.Client
}

func NewMilvusRepository(ctx context.Context, cfg *config.Config) (vector.VectorDB, client.Client, error) {
	milvus, err := client.NewGrpcClient(ctx, cfg.Milvus.Host)
	if err != nil {
		return nil, nil, err
	}

	return &Repository{milvus: milvus}, milvus, nil
}

func (r Repository) GetTopK(ctx context.Context, orgID string, k int, search []float32) ([]client.SearchResult, error) {
	orgID = "_" + strings.ReplaceAll(orgID, "-", "")
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
		ctx,              // ctx
		orgID,            // CollectionName
		[]string{},       // partitionNames
		"",               // expr
		[]string{"text"}, // outputFields
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

func (r Repository) SaveDoc(ctx context.Context, orgID, docID string, chunks []string, embeddings [][][]float32) error {
	orgID = "_" + strings.ReplaceAll(orgID, "-", "")
	ids := make([]string, len(embeddings))
	for i, _ := range ids {
		ids[i] = docID
	}
	idColumn := entity.NewColumnVarChar("id", ids)
	textColumn := entity.NewColumnVarChar("text", chunks)
	embeddingColumn := entity.NewColumnFloatVector("embedding", 3072, bind(embeddings))

	ok, err := r.milvus.HasCollection(ctx, orgID)
	if err != nil {
		return err
	}
	if !ok {
		schema := r.createSchema(orgID)
		err := r.milvus.CreateCollection(
			ctx, // ctx
			schema,
			2, // shardNum
		)
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
			ctx,         // ctx
			orgID,       // CollectionName
			"embedding", // fieldName
			idx,         // entity.Index
			false,       // async
		)
		if err != nil {
			return err
		}
	}
	_, err = r.milvus.Insert(
		ctx,      // ctx
		orgID,    // CollectionName
		"",       // partitionName
		idColumn, // columnarData
		textColumn,
		embeddingColumn, // columnarData
	)
	if err != nil {
		return err
	}
	err = r.milvus.LoadCollection(
		ctx,   // ctx
		orgID, // CollectionName
		false, // async
	)
	return nil
}

func bind(embeddings [][][]float32) [][]float32 {
	result := make([][]float32, 0)
	for _, embedding := range embeddings {
		if len(embedding) == 0 {
			result = append(result, emptyDym())
		} else {
			result = append(result, embedding[0])
		}
	}
	return result
}

func emptyDym() []float32 {
	var result []float32
	for i := 0; i < 3072; i++ {
		result = append(result, 0)
	}
	return result
}

func (r Repository) DeleteDoc(ctx context.Context, orgID string, id string) error {
	orgID = "_" + strings.ReplaceAll(orgID, "-", "")
	expr := fmt.Sprintf("id == '%s'", strings.Trim(id, "/"))
	err := r.milvus.Delete(
		ctx,   // ctx
		orgID, // collection name
		"",    // partition name
		expr,  // expr
	)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) createSchema(orgID string) *entity.Schema {
	return &entity.Schema{
		CollectionName: orgID,
		Description:    "testing",
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				AutoID:     false,
				TypeParams: map[string]string{
					"max_length": "10000",
				},
			},
			{
				Name:       "text",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: false,
				AutoID:     false,
				TypeParams: map[string]string{
					"max_length": "10000",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "3072",
				},
			},
		},
		EnableDynamicField: true,
	}
}
