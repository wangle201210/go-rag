package websearch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/retriever"
)

// newTestServer 启动一个可控的 you.com 假服务，返回 server 及最近一次收到的请求信息
func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func TestRetrieve_RequestConstruction(t *testing.T) {
	var gotPath, gotAPIKey, gotQuery, gotCount string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAPIKey = r.Header.Get("X-API-Key")
		gotQuery = r.URL.Query().Get("query")
		gotCount = r.URL.Query().Get("count")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[]}}`))
	})

	r, err := NewRetriever(&Config{APIKey: "test-key", BaseURL: srv.URL + "/v1/search", Count: 3})
	if err != nil {
		t.Fatalf("NewRetriever failed: %v", err)
	}
	if _, err = r.Retrieve(context.Background(), "hello world"); err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}

	if gotPath != "/v1/search" {
		t.Errorf("path = %q, want /v1/search", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Errorf("X-API-Key header = %q, want test-key", gotAPIKey)
	}
	if gotQuery != "hello world" {
		t.Errorf("query param = %q, want %q", gotQuery, "hello world")
	}
	if gotCount != "3" {
		t.Errorf("count param = %q, want 3", gotCount)
	}
}

func TestRetrieve_MapsWebAndNews(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{
			"web":[{"url":"https://a.com","title":"A","description":"descA","snippets":["s1","s2"]}],
			"news":[{"url":"https://b.com","title":"B","description":"descB","snippets":["s3"]}]
		}}`))
	})

	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	docs, err := r.Retrieve(context.Background(), "q")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("got %d docs, want 2", len(docs))
	}
	if docs[0].ID != "https://a.com" || docs[0].MetaData["url"] != "https://a.com" {
		t.Errorf("web doc ID/url mismatch: %+v", docs[0])
	}
	if docs[1].ID != "https://b.com" {
		t.Errorf("news doc ID mismatch: %+v", docs[1])
	}
	if docs[1].MetaData["source"] != sourceName {
		t.Errorf("source metadata = %v, want %v", docs[1].MetaData["source"], sourceName)
	}
}

func TestRetrieve_NewsAbsent(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[{"url":"https://a.com","title":"A","description":"d"}]}}`))
	})

	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	docs, err := r.Retrieve(context.Background(), "q")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("got %d docs, want 1 (news absent should not break mapping)", len(docs))
	}
}

func TestRetrieve_OptionalFieldsMissing(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 只包含 url，title/description/snippets 均缺失
		w.Write([]byte(`{"results":{"web":[{"url":"https://a.com"}]}}`))
	})

	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	docs, err := r.Retrieve(context.Background(), "q")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("got %d docs, want 1", len(docs))
	}
	if docs[0].ID != "https://a.com" {
		t.Errorf("ID = %q, want https://a.com", docs[0].ID)
	}
	if docs[0].Content != "" {
		t.Errorf("Content = %q, want empty when title/description/snippets absent", docs[0].Content)
	}
}

func TestRetrieve_SnippetsJoined(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[{"url":"https://a.com","title":"T","description":"D","snippets":["s1","s2","s3"]}]}}`))
	})

	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	docs, err := r.Retrieve(context.Background(), "q")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	want := "T\nD\ns1\ns2\ns3"
	if docs[0].Content != want {
		t.Errorf("Content = %q, want %q", docs[0].Content, want)
	}
}

func TestRetrieve_ScoreDescendingOrder(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[
			{"url":"https://a.com","title":"A"},
			{"url":"https://b.com","title":"B"},
			{"url":"https://c.com","title":"C"}
		]}}`))
	})

	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	docs, err := r.Retrieve(context.Background(), "q")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if len(docs) != 3 {
		t.Fatalf("got %d docs, want 3", len(docs))
	}
	for i := 1; i < len(docs); i++ {
		if docs[i].Score() >= docs[i-1].Score() {
			t.Errorf("scores not strictly descending at %d: %v >= %v", i, docs[i].Score(), docs[i-1].Score())
		}
	}
}

func TestRetrieve_ErrorStatusCodes(t *testing.T) {
	cases := []struct {
		status int
	}{
		{http.StatusUnauthorized},
		{http.StatusForbidden},
		{http.StatusUnprocessableEntity},
		{http.StatusTooManyRequests},
		{http.StatusInternalServerError},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.status), func(t *testing.T) {
			srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(`{"error":"boom"}`))
			})
			r, _ := NewRetriever(&Config{APIKey: "super-secret-key", BaseURL: srv.URL})
			_, err := r.Retrieve(context.Background(), "q")
			if err == nil {
				t.Fatalf("expected error for status %d, got nil", tc.status)
			}
			if got := err.Error(); containsKey(got, "super-secret-key") {
				t.Errorf("error message leaked API key: %q", got)
			}
		})
	}
}

func TestRetrieve_MalformedJSON(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not json`))
	})
	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	_, err := r.Retrieve(context.Background(), "q")
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestRetrieve_Timeout(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[]}}`))
	})
	r, _ := NewRetriever(&Config{APIKey: "super-secret-key", BaseURL: srv.URL})
	r.client.Timeout = 20 * time.Millisecond // 覆盖为极短超时以加速测试

	_, err := r.Retrieve(context.Background(), "q")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if containsKey(err.Error(), "super-secret-key") {
		t.Errorf("timeout error leaked API key: %q", err.Error())
	}
}

func TestRetrieve_NoAPIKeyAnywhere(t *testing.T) {
	t.Setenv("YDC_API_KEY", "")
	r, _ := NewRetriever(&Config{})
	_, err := r.Retrieve(context.Background(), "q")
	if err == nil {
		t.Fatal("expected error when no API key is configured anywhere")
	}
}

func TestRetrieve_CountDefaultAndCap(t *testing.T) {
	var gotCount string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotCount = r.URL.Query().Get("count")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[]}}`))
	})

	t.Run("default", func(t *testing.T) {
		r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
		if _, err := r.Retrieve(context.Background(), "q"); err != nil {
			t.Fatalf("Retrieve failed: %v", err)
		}
		if gotCount != strconv.Itoa(defaultCount) {
			t.Errorf("count = %q, want default %d", gotCount, defaultCount)
		}
	})

	t.Run("cap at construction", func(t *testing.T) {
		r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL, Count: 999})
		if _, err := r.Retrieve(context.Background(), "q"); err != nil {
			t.Fatalf("Retrieve failed: %v", err)
		}
		if gotCount != strconv.Itoa(maxCount) {
			t.Errorf("count = %q, want capped %d", gotCount, maxCount)
		}
	})

	t.Run("cap via WithTopK option", func(t *testing.T) {
		r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
		if _, err := r.Retrieve(context.Background(), "q", retriever.WithTopK(50)); err != nil {
			t.Fatalf("Retrieve failed: %v", err)
		}
		if gotCount != strconv.Itoa(maxCount) {
			t.Errorf("count = %q, want capped %d", gotCount, maxCount)
		}
	})
}

func TestRetrieve_APIKeyPrecedence(t *testing.T) {
	var gotKey string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-API-Key")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"web":[]}}`))
	})

	t.Run("yaml key wins over env", func(t *testing.T) {
		t.Setenv("YDC_API_KEY", "env-key")
		r, _ := NewRetriever(&Config{APIKey: "yaml-key", BaseURL: srv.URL})
		if _, err := r.Retrieve(context.Background(), "q"); err != nil {
			t.Fatalf("Retrieve failed: %v", err)
		}
		if gotKey != "yaml-key" {
			t.Errorf("key used = %q, want yaml-key", gotKey)
		}
	})

	t.Run("falls back to env when yaml empty", func(t *testing.T) {
		t.Setenv("YDC_API_KEY", "env-key")
		r, _ := NewRetriever(&Config{BaseURL: srv.URL})
		if _, err := r.Retrieve(context.Background(), "q"); err != nil {
			t.Fatalf("Retrieve failed: %v", err)
		}
		if gotKey != "env-key" {
			t.Errorf("key used = %q, want env-key", gotKey)
		}
	})
}

func containsKey(s, key string) bool {
	return len(key) > 0 && strings.Contains(s, key)
}

func TestRetrieve_TopKTruncatesTotalResults(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{
			"web":[
				{"url":"https://a.com","title":"A","description":"d"},
				{"url":"https://b.com","title":"B","description":"d"}
			],
			"news":[
				{"url":"https://c.com","title":"C","description":"d"},
				{"url":"https://d.com","title":"D","description":"d"}
			]
		}}`))
	})

	r, _ := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	docs, err := r.Retrieve(context.Background(), "q", retriever.WithTopK(2))
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	// web 与 news 双通道合计可能超过 TopK，Retrieve 需按总数截断，
	// 且 web 结果（分数更高）优先保留
	if len(docs) != 2 {
		t.Fatalf("got %d docs, want 2 (TopK must bound total results)", len(docs))
	}
	if docs[0].ID != "https://a.com" || docs[1].ID != "https://b.com" {
		t.Errorf("expected web docs kept first after truncation, got %v, %v", docs[0].ID, docs[1].ID)
	}
}

func TestRetrieve_ZeroTopKUsesConfiguredCount(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, req *http.Request) {
		if got := req.URL.Query().Get("count"); got != "3" {
			t.Errorf("count = %q, want configured count 3", got)
		}
		w.Write([]byte(`{"results":{"web":[]}}`))
	})
	r, err := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL, Count: 3})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Retrieve(context.Background(), "q", retriever.WithTopK(0)); err != nil {
		t.Fatal(err)
	}
}

func TestRetrieve_RejectsInvalidParametersBeforeRequest(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, req *http.Request) {
		t.Error("invalid parameters must not reach the provider")
	})
	r, err := NewRetriever(&Config{APIKey: "k", BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{"", " \n\t "} {
		if _, err := r.Retrieve(context.Background(), query); err == nil {
			t.Error("expected blank query error")
		}
	}
	if _, err := r.Retrieve(context.Background(), "q", retriever.WithTopK(-1)); err == nil {
		t.Error("expected negative top_k error")
	}
}
