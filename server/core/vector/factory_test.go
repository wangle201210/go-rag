package vector

import (
	"testing"
)

func TestNewVectorStore_ES(t *testing.T) {
	cfg := &Config{
		Type:      "es",
		IndexName: "test-index",
		ES: &ESConfig{
			Address:  "http://localhost:9200",
			Username: "",
			Password: "",
		},
	}

	store, err := NewVectorStore(cfg)
	if err != nil {
		t.Fatalf("Failed to create ES vector store: %v", err)
	}

	if store == nil {
		t.Fatal("Expected non-nil vector store")
	}

	if _, ok := store.(*ESVectorStore); !ok {
		t.Fatal("Expected ESVectorStore type")
	}
}

func TestNewVectorStore_Qdrant(t *testing.T) {
	cfg := &Config{
		Type:      "qdrant",
		IndexName: "test-index",
		Qdrant: &QdrantConfig{
			Address: "http://localhost:6333",
			APIKey:  "",
		},
	}

	store, err := NewVectorStore(cfg)
	if err != nil {
		t.Fatalf("Failed to create Qdrant vector store: %v", err)
	}

	if store == nil {
		t.Fatal("Expected non-nil vector store")
	}

	if _, ok := store.(*QdrantVectorStore); !ok {
		t.Fatal("Expected QdrantVectorStore type")
	}
}

func TestNewVectorStore_InvalidType(t *testing.T) {
	cfg := &Config{
		Type:      "invalid",
		IndexName: "test-index",
	}

	_, err := NewVectorStore(cfg)
	if err == nil {
		t.Fatal("Expected error for invalid type")
	}
}

func TestNewVectorStore_MissingConfig(t *testing.T) {
	cfg := &Config{
		Type:      "es",
		IndexName: "test-index",
		// ES config is nil
	}

	_, err := NewVectorStore(cfg)
	if err == nil {
		t.Fatal("Expected error for missing ES config")
	}
}
