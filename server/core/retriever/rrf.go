package retriever

import (
	"sort"

	"github.com/cloudwego/eino/schema"
)

// RRFFusion 对多路召回结果执行 RRF（Reciprocal Rank Fusion）融合。
// inputs 每个元素是一路召回的有序文档列表（按相关性降序）。
func RRFFusion(inputs [][]*schema.Document) []*schema.Document {
	const k = 60

	docScores := make(map[string]float64)
	docMap := make(map[string]*schema.Document)

	for _, docs := range inputs {
		for rank, doc := range docs {
			if doc.ID == "" {
				continue
			}
			docScores[doc.ID] += 1.0 / float64(k+rank+1)
			if _, exists := docMap[doc.ID]; !exists {
				docMap[doc.ID] = doc
			}
		}
	}

	result := make([]*schema.Document, 0, len(docMap))
	for id, score := range docScores {
		doc := docMap[id]
		doc.WithScore(score)
		result = append(result, doc)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score() > result[j].Score()
	})

	return result
}
