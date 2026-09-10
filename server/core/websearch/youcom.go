// Package websearch 基于 You.com Search API 实现联网检索能力，
// 作为 eino retriever.Retriever 的一种实现，独立于现有的向量检索链路。
package websearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	// defaultBaseURL You.com Search API 默认地址
	defaultBaseURL = "https://ydc-index.io/v1/search"
	// defaultCount 默认返回结果数量
	defaultCount = 5
	// maxCount 本组件限制的最大返回结果数量
	maxCount = 20
	// defaultTimeout 默认请求超时时间
	defaultTimeout = 10 * time.Second
	// sourceName 结果来源标识，写入 MetaData
	sourceName = "you.com"
	// apiKeyEnv You.com API Key 的环境变量名
	apiKeyEnv = "YDC_API_KEY"
)

// Config You.com 联网检索配置
type Config struct {
	APIKey     string // You.com Search API Key，留空时读取环境变量 YDC_API_KEY
	BaseURL    string // 可选，默认 https://ydc-index.io/v1/search
	Count      int    // 默认返回结果数量，默认 5，最大 20
	TimeoutSec int    // 请求超时时间（秒），默认 10
}

// YouComRetriever 基于 You.com Search API 实现的 eino Retriever
type YouComRetriever struct {
	cfg    *Config
	client *http.Client
}

// youComResponse You.com Search API 响应结构（仅保留用得到的字段）
type youComResponse struct {
	Results struct {
		Web  []youComResult `json:"web"`
		News []youComResult `json:"news"`
	} `json:"results"`
}

// youComResult 单条搜索结果，除 url/title/description/snippets 外的字段均视为可选
type youComResult struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Snippets    []string `json:"snippets"`
}

// NewRetriever 创建 You.com 联网检索 retriever
func NewRetriever(cfg *Config) (*YouComRetriever, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	c := *cfg // 拷贝一份，避免修改调用方传入的配置对象

	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.Count <= 0 {
		c.Count = defaultCount
	}
	if c.Count > maxCount {
		c.Count = maxCount
	}

	timeout := defaultTimeout
	if c.TimeoutSec > 0 {
		timeout = time.Duration(c.TimeoutSec) * time.Second
	}

	return &YouComRetriever{
		cfg: &c,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// resolveAPIKey 解析 API Key：优先使用配置文件中的 apiKey，否则读取环境变量 YDC_API_KEY
func (r *YouComRetriever) resolveAPIKey() string {
	if r.cfg.APIKey != "" {
		return r.cfg.APIKey
	}
	return os.Getenv(apiKeyEnv)
}

// Retrieve 调用 You.com Search API 检索网页结果，实现 eino retriever.Retriever 接口
func (r *YouComRetriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	if strings.TrimSpace(query) == "" {
		return nil, gerror.New("联网检索问题不能为空")
	}
	key := r.resolveAPIKey()
	if key == "" {
		return nil, fmt.Errorf("未配置 You.com API Key，请设置环境变量 %s 或在配置文件 websearch.youcom.apiKey 中填写", apiKeyEnv)
	}

	options := &retriever.Options{}
	retriever.GetCommonOptions(options, opts...)

	if options.TopK != nil && *options.TopK < 0 {
		return nil, gerror.New("联网检索 top_k 不能为负数")
	}
	count := r.cfg.Count
	if options.TopK != nil && *options.TopK > 0 {
		count = *options.TopK
	}
	if count > maxCount {
		count = maxCount
	}

	req, err := r.buildRequest(ctx, key, query, count)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, mapRequestError(ctx, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 you.com 响应失败: %w", err)
	}

	if statusErr := mapStatusError(resp.StatusCode); statusErr != nil {
		return nil, statusErr
	}

	var payload youComResponse
	if err = json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("解析 you.com 响应失败: %w", err)
	}

	docs := toDocuments(&payload)
	// web 与 news 双通道各返回最多 count 条，这里按总数截断，
	// 保证调用方拿到的结果数量与 TopK 语义一致（与 /v1/retriever 行为对齐）
	if len(docs) > count {
		docs = docs[:count]
	}
	return docs, nil
}

// buildRequest 构建 you.com search 请求：GET + X-API-Key header + query/count 参数
func (r *YouComRetriever) buildRequest(ctx context.Context, key, query string, count int) (*http.Request, error) {
	u, err := url.Parse(r.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("you.com baseURL 非法: %w", err)
	}
	q := u.Query()
	q.Set("query", query)
	q.Set("count", strconv.Itoa(count))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("构建 you.com 请求失败: %w", err)
	}
	req.Header.Set("X-API-Key", key)
	return req, nil
}

// mapRequestError 将底层网络错误转换为更明确的提示，never 回显 API Key
func mapRequestError(ctx context.Context, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("you.com 请求超时: %w", err)
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("you.com 请求超时: %w", err)
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return fmt.Errorf("you.com 请求已取消: %w", ctxErr)
	}
	return fmt.Errorf("you.com 请求失败: %w", err)
}

// mapStatusError 将 you.com 的错误状态码映射为可读的错误信息，不回显 API Key
func mapStatusError(statusCode int) error {
	switch {
	case statusCode == http.StatusOK:
		return nil
	case statusCode == http.StatusUnauthorized:
		return fmt.Errorf("you.com 鉴权失败(401)，请检查 API Key 是否正确")
	case statusCode == http.StatusForbidden:
		return fmt.Errorf("you.com 请求被拒绝(403)，请检查请求地址或账号权限")
	case statusCode == http.StatusUnprocessableEntity:
		return fmt.Errorf("you.com 请求参数非法(422)，请检查检索参数")
	case statusCode == http.StatusTooManyRequests:
		return fmt.Errorf("you.com 请求超出速率限制(429)，请稍后重试")
	case statusCode >= http.StatusInternalServerError:
		return fmt.Errorf("you.com 服务端错误(%d)，可稍后重试", statusCode)
	default:
		return fmt.Errorf("you.com 返回非预期状态码(%d)", statusCode)
	}
}

// toDocuments 将 you.com 的搜索结果转换为 schema.Document 列表
// results.web 与 results.news 都会被收录，results.news 缺失时按空处理
func toDocuments(payload *youComResponse) []*schema.Document {
	results := make([]youComResult, 0, len(payload.Results.Web)+len(payload.Results.News))
	results = append(results, payload.Results.Web...)
	results = append(results, payload.Results.News...)

	docs := make([]*schema.Document, 0, len(results))
	for i, item := range results {
		doc := &schema.Document{
			ID:      item.URL,
			Content: buildContent(item),
			MetaData: map[string]any{
				"url":    item.URL,
				"title":  item.Title,
				"source": sourceName,
			},
		}
		// you.com 按相关性排好序返回，没有分数，这里用排名倒序模拟一个递减分数
		doc.WithScore(rankScore(i))
		docs = append(docs, doc)
	}
	return docs
}

// buildContent 拼接 title + description + snippets 作为文档正文
func buildContent(item youComResult) string {
	parts := make([]string, 0, 2+len(item.Snippets))
	if item.Title != "" {
		parts = append(parts, item.Title)
	}
	if item.Description != "" {
		parts = append(parts, item.Description)
	}
	if len(item.Snippets) > 0 {
		parts = append(parts, strings.Join(item.Snippets, "\n"))
	}
	return strings.Join(parts, "\n")
}

// rankScore 依据排名生成一个递减的分数，排名越靠前分数越高
func rankScore(rank int) float64 {
	const epsilon = 0.001
	return 1.0 - float64(rank)*epsilon
}
