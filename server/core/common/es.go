package common

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// createIndex create index for example in add_documents.go.
func createIndex(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	_, err := create.NewCreateFunc(client)(indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				FieldContent:  types.NewTextProperty(),
				FieldExtra:    types.NewTextProperty(),
				KnowledgeName: types.NewKeywordProperty(),
				FieldContentVector: &types.DenseVectorProperty{
					Dims:       Of(1024), // same as embedding dimensions
					Index:      Of(true),
					Similarity: Of("cosine"),
				},
			},
		},
	}).Do(ctx)

	return err
}

func CreateIndexIfNotExists(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	indexExists, err := exists.NewExistsFunc(client)(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if indexExists {
		return nil
	}
	err = createIndex(ctx, client, indexName)
	return err
}
