package retriever

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

func TestBM25RequestAndResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Size  int `json:"size"`
			Query struct {
				Bool struct {
					Must []struct {
						Match map[string]struct {
							Query string `json:"query"`
						} `json:"match"`
					} `json:"must"`
					Filter []struct {
						Term map[string]struct {
							Value string `json:"value"`
						} `json:"term"`
					} `json:"filter"`
					MustNot []struct {
						Terms map[string][]string `json:"terms"`
					} `json:"must_not"`
				} `json:"bool"`
			} `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		b := body.Query.Bool
		if r.URL.Path != "/chunks/_search" || body.Size != 3 {
			t.Errorf("wrong index/limit: %s %+v", r.URL.Path, body)
		}
		if len(b.Must) != 1 || b.Must[0].Match[coretypes.FieldContent].Query != "search words" {
			t.Errorf("missing content query: %+v", b.Must)
		}
		if len(b.Filter) != 1 || b.Filter[0].Term[coretypes.KnowledgeName].Value != "private-kb" {
			t.Errorf("missing exact knowledge filter: %+v", b.Filter)
		}
		if len(b.MustNot) != 1 || len(b.MustNot[0].Terms["_id"]) != 1 || b.MustNot[0].Terms["_id"][0] != "excluded" {
			t.Errorf("missing exclusions: %+v", b.MustNot)
		}
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"hits":{"hits":[{"_id":"one","_score":2.5,"_source":{"content":"found","_knowledge_name":"private-kb"}}]}}`)
	}))
	defer srv.Close()
	client, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{srv.URL}, DisableRetry: true})
	if err != nil {
		t.Fatal(err)
	}
	docs, err := Bm25Retrieve(context.Background(), client, "chunks", "private-kb", "search words", []string{"excluded"}, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 1 || docs[0].ID != "one" || docs[0].Content != "found" || docs[0].Score() != 2.5 {
		t.Fatalf("unexpected docs: %+v", docs)
	}
}

func TestBM25Failure(t *testing.T) {
	for _, status := range []int{http.StatusServiceUnavailable, http.StatusOK} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(status)
				if status == http.StatusOK {
					fmt.Fprint(w, `{"hits":{"hits":[{"_id":"broken","_source":{"unexpected":true}}]}}`)
				} else {
					fmt.Fprint(w, `{"error":{"type":"unavailable","reason":"offline"},"status":503}`)
				}
			}))
			defer srv.Close()
			client, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{srv.URL}, DisableRetry: true})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Bm25Retrieve(context.Background(), client, "chunks", "kb", "q", nil, 3); err == nil {
				t.Fatal("expected search/parser error")
			}
		})
	}
	if _, err := Bm25Retrieve(context.Background(), nil, "chunks", "kb", "q", nil, 3); err == nil {
		t.Fatal("expected missing client error")
	}
}
