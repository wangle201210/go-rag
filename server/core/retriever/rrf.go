package retriever

import (
	"sort"

	"github.com/cloudwego/eino/schema"
)

// RRFFusion 融合按相关性降序排列的召回列表，每个文档在每一路最多贡献一次分数。
// 同分时保留首次出现顺序，不修改调用方的文档或元数据。
func RRFFusion(inputs [][]*schema.Document) []*schema.Document {
	const k = 60
	var (
		docScores = make(map[string]float64)
		docMap    = make(map[string]*schema.Document)
		result    = make([]*schema.Document, 0)
	)
	for _, docs := range inputs {
		seen := make(map[string]bool)
		rank := 0
		for _, doc := range docs {
			if doc == nil || doc.ID == "" || seen[doc.ID] {
				continue
			}
			seen[doc.ID] = true
			rank++
			docScores[doc.ID] += 1.0 / float64(k+rank)
			if _, exists := docMap[doc.ID]; !exists {
				cloned := *doc
				cloned.MetaData = make(map[string]any, len(doc.MetaData)+1)
				for key, value := range doc.MetaData {
					cloned.MetaData[key] = value
				}
				docMap[doc.ID] = &cloned
				result = append(result, &cloned)
			}
		}
	}
	for _, doc := range result {
		doc.WithScore(docScores[doc.ID])
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score() > result[j].Score()
	})
	return result
}
