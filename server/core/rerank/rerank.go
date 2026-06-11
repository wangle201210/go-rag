package rerank

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

type Conf struct {
	Model           string `json:"model"`
	ReturnDocuments bool   `json:"return_documents"`
	MaxChunksPerDoc int    `json:"max_chunks_per_doc"`
	OverlapTokens   int    `json:"overlap_tokens"`
	url             string
	apiKey          string
}
type Data struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n"`
}

type Req struct {
	*Data
	*Conf
}

type Result struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
}

type Resp struct {
	ID      string    `json:"id"`
	Results []*Result `json:"results"`
}

var rerankCfg *Conf

func NewRerank(ctx context.Context, query string, docs []*schema.Document, topK int) (output []*schema.Document, err error) {
	output, err = rerank(ctx, query, docs, topK)
	if err != nil {
		g.Log().Warningf(ctx, "Rerank remote failed, fallback to RRF: %v", err)
		return rrfFallback(docs, topK), nil
	}
	return output, nil
}

// rrfFallback 用 RRF 算法对单路结果（已按原始 score 排好序）重新打分并截断到 topK。
// 单路 RRF 本质上保持原顺序，但把分数统一到可比较的 1/(k+rank) 量纲。
func rrfFallback(docs []*schema.Document, topK int) []*schema.Document {
	const k = 60
	for rank, doc := range docs {
		doc.WithScore(1.0 / float64(k+rank+1))
	}
	if topK > 0 && len(docs) > topK {
		return docs[:topK]
	}
	return docs
}

func GetConf(ctx context.Context) *Conf {
	if rerankCfg != nil {
		return rerankCfg
	}
	baseUrl := g.Cfg().MustGet(ctx, "rerank.baseURL").String()
	apiKey := g.Cfg().MustGet(ctx, "rerank.apiKey").String()
	model := g.Cfg().MustGet(ctx, "rerank.model").String()
	url := fmt.Sprintf("%s/rerank", baseUrl)
	rerankCfg = &Conf{
		apiKey:          apiKey,
		Model:           model,
		ReturnDocuments: false,
		MaxChunksPerDoc: 1024,
		OverlapTokens:   80,
		url:             url,
	}
	return rerankCfg
}

func rerank(ctx context.Context, query string, docs []*schema.Document, topK int) (output []*schema.Document, err error) {
	data := &Data{
		Query: query,
		TopN:  topK,
	}
	// g.Log().Infof(ctx, "docs num: %d", len(docs))
	for _, doc := range docs {
		data.Documents = append(data.Documents, doc.Content)
	}
	// 重排
	results, err := rerankDoHttp(ctx, data)
	if err != nil {
		return
	}
	// 重新组装数据
	for _, result := range results {
		doc := docs[result.Index]
		// g.Log().Infof(ctx, "content: %s, score_old: %f, score_new: %f", doc.Content, doc.Score(), result.RelevanceScore)
		doc.WithScore(result.RelevanceScore)
		output = append(output, docs[result.Index])
	}
	return
}

func rerankDoHttp(ctx context.Context, data *Data) ([]*Result, error) {
	cfg := GetConf(ctx)
	reqData := &Req{
		Data: data,
		Conf: cfg,
	}

	marshal, err := sonic.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewReader(marshal)
	request, err := http.NewRequest("POST", cfg.url, payload)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.apiKey))
	do, err := g.Client().Do(request)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	body, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	res := Resp{}
	err = sonic.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res.Results, nil
}
