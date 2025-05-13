package embedding

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/embedding"
)

type Sparse struct {
	client *http.Client
	url    string
}

type sparseRequest struct {
	Texts        []string `json:"texts"`
	ReturnDense  *bool    `json:"return_dense,omitempty"`
	ReturnSparse *bool    `json:"return_sparse,omitempty"`
}

func NewSparse() *Sparse {
	return &Sparse{
		client: http.DefaultClient,
		url:    "http://localhost:8082/encode",
	}
}

type sparseResponse struct {
	SparseEmbeddings []*sparseData `json:"sparse_embeddings"`
}

type sparseData struct {
	Raw      map[string]float64 `json:"raw"`
	Readable map[string]float64 `json:"readable"`
}

func (x *Sparse) EmbedStrings(ctx context.Context, texts []string, _ ...embedding.Option) (res []map[string]float64, err error) {
	sr := &sparseRequest{
		Texts: texts,
	}
	marshal, _ := sonic.Marshal(sr)
	req, _ := http.NewRequest("POST", x.url, bytes.NewReader(marshal))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	response, _ := x.client.Do(req)

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	resData := &sparseResponse{}
	err = sonic.Unmarshal(body, &resData)
	if err != nil {
		return nil, err
	}
	res = make([]map[string]float64, 0, len(resData.SparseEmbeddings))
	for _, se := range resData.SparseEmbeddings {
		res = append(res, se.Raw)
	}
	return res, nil
}
