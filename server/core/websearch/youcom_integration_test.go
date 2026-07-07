package websearch

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestYouComRetrieverIntegration 是唯一一个会真正访问 You.com 线上接口的测试，
// 未设置 YDC_API_KEY 时自动跳过，避免 CI / 无 key 环境失败。
func TestYouComRetrieverIntegration(t *testing.T) {
	if os.Getenv("YDC_API_KEY") == "" {
		t.Skip("YDC_API_KEY 未设置，跳过 you.com 线上集成测试")
	}

	r, err := NewRetriever(&Config{Count: 3, TimeoutSec: 10})
	if err != nil {
		t.Fatalf("NewRetriever failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	docs, err := r.Retrieve(ctx, "You.com API")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if len(docs) == 0 {
		t.Fatal("expected at least 1 document from live you.com search")
	}
	if docs[0].Content == "" {
		t.Error("expected non-empty Content in first document")
	}
	if docs[0].MetaData["url"] == nil || docs[0].MetaData["url"] == "" {
		t.Error("expected non-empty MetaData[url] in first document")
	}
}
