package retriever

import (
	"math"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestRRFFusion(t *testing.T) {
	a := (&schema.Document{ID: "a", Content: "first", MetaData: map[string]any{"source": "original"}}).WithScore(0.9)
	b := &schema.Document{ID: "b"}
	c := &schema.Document{ID: "c"}
	got := RRFFusion([][]*schema.Document{
		{nil, {}, a, a, b},
		{b, c},
	})
	if len(got) != 3 || got[0].ID != "b" || got[1].ID != "a" || got[2].ID != "c" {
		t.Fatalf("unexpected fusion order: %+v", got)
	}
	if math.Abs(got[0].Score()-(1.0/62+1.0/61)) > 1e-12 || math.Abs(got[1].Score()-1.0/61) > 1e-12 {
		t.Fatalf("unexpected scores: %v, %v", got[0].Score(), got[1].Score())
	}
	if a.Score() != 0.9 || got[1] == a || got[1].Content != "first" {
		t.Fatal("fusion mutated the input document or lost its content")
	}
	got[1].MetaData["source"] = "changed"
	if a.MetaData["source"] != "original" {
		t.Fatal("fusion shares mutable metadata with input")
	}
}

func TestRRFFusionStableTies(t *testing.T) {
	for i := 0; i < 100; i++ {
		got := RRFFusion([][]*schema.Document{{{ID: "b"}}, {{ID: "a"}}})
		if len(got) != 2 || got[0].ID != "b" || got[1].ID != "a" {
			t.Fatal("ties must preserve first appearance")
		}
	}
	if got := RRFFusion(nil); len(got) != 0 {
		t.Fatalf("expected empty fusion, got %v", got)
	}
}
