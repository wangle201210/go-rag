package common

import (
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestRemoveDuplicates(t *testing.T) {
	docs := []*schema.Document{
		{
			ID: "1",
			MetaData: map[string]any{
				"foo": "bar",
			},
		},
		{
			ID: "2",
		},
		{
			ID: "3",
		},
		{
			ID: "1",
		},
	}

	docs = RemoveDuplicates(docs, func(t *schema.Document) string {
		return t.ID
	})
	for i, doc := range docs {
		t.Logf("i: %d, doc_id: %v", i, doc.ID)
	}
}
