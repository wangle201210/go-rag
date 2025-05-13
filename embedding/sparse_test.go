package embedding

import (
	"context"
	"testing"
)

func TestNewSparse(t *testing.T) {
	sp := NewSparse()
	res, err := sp.EmbedStrings(context.Background(), []string{"hello"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
