package common

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func GetESClient(ctx context.Context) (client *elasticsearch.Client, err error) {
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return nil, err
	}
	if err = createIndexIfNotExists(ctx, client); err != nil {
		log.Printf("createIndex of es8 failed, err=%v", err)
		return nil, err
	}
	return
}

// createIndex create index for example in add_documents.go.
func createIndex(ctx context.Context, client *elasticsearch.Client) error {
	_, err := create.NewCreateFunc(client)(IndexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				FieldContent: types.NewTextProperty(),
				FieldExtra:   types.NewTextProperty(),
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

func createIndexIfNotExists(ctx context.Context, client *elasticsearch.Client) error {
	indexExists, err := exists.NewExistsFunc(client)(IndexName).Do(ctx)
	if err != nil {
		return err
	}
	if indexExists {
		return nil
	}
	err = createIndex(ctx, client)
	return err
}
