package indexer

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// docAddIDAndMerge component initialization function of node 'Lambda1' in graph 't'
func docAddIDAndMerge(ctx context.Context, docs []*schema.Document) (output []*schema.Document, err error) {
	if len(docs) == 0 {
		return docs, nil
	}
	for _, doc := range docs {
		doc.ID = uuid.New().String() // 覆盖之前的id
	}
	switch docs[0].MetaData[file.MetaKeyExtension] {
	case ".md":
		return mergeMD(ctx, docs)
	case ".xlsx":
		return mergeXLSX(ctx, docs)
	default:
		return docs, nil
	}
}

func mergeMD(ctx context.Context, docs []*schema.Document) (output []*schema.Document, err error) {
	ndocs := make([]*schema.Document, 0, len(docs))
	var nd *schema.Document
	maxLen := 512
	for _, doc := range docs {
		// 不是同一个文件的就不要放一起了
		if nd != nil && doc.MetaData[file.MetaKeySource] != nd.MetaData[file.MetaKeySource] {
			ndocs = append(ndocs, nd)
			nd = nil
		}
		// 两个文档长度之和大于maxLen就不要放一起了
		if nd != nil && len(nd.Content)+len(doc.Content) > maxLen {
			ndocs = append(ndocs, nd)
			nd = nil
		}
		// 不是同一个一级标题的就不要放一起了
		if nd != nil && doc.MetaData[coretypes.Title1] != nd.MetaData[coretypes.Title1] {
			ndocs = append(ndocs, nd)
			nd = nil
		}
		// 不是同一个二级标题的就不要放一起了
		// 如果nd的h2是nil，证明之前只有h1,且两个的h1相等，则直接合并
		if nd != nil && nd.MetaData[coretypes.Title2] != nil && doc.MetaData[coretypes.Title2] != nd.MetaData[coretypes.Title2] {
			ndocs = append(ndocs, nd)
			nd = nil
		}
		if nd == nil {
			nd = doc
		} else {
			mergeTitle(nd, doc, coretypes.Title2)
			mergeTitle(nd, doc, coretypes.Title3)
			nd.Content += doc.Content
		}
	}
	if nd != nil {
		ndocs = append(ndocs, nd)
	}
	for _, ndoc := range ndocs {
		ndoc.Content = getMdContentWithTitle(ndoc)
	}
	return ndocs, nil
}

func mergeXLSX(ctx context.Context, docs []*schema.Document) (output []*schema.Document, err error) {
	for _, doc := range docs {
		marshal, _ := sonic.Marshal(doc.MetaData[coretypes.XlsxRow])
		doc.Content = string(marshal)
	}
	return docs, nil
}

func mergeTitle(orgDoc, addDoc *schema.Document, key string) {
	// 相等就不管了
	if orgDoc.MetaData[key] == addDoc.MetaData[key] {
		return
	}
	var title []string
	if orgDoc.MetaData[key] != nil {
		title = append(title, orgDoc.MetaData[key].(string))
	}
	if addDoc.MetaData[key] != nil {
		title = append(title, addDoc.MetaData[key].(string))
	}
	if len(title) > 0 {
		orgDoc.MetaData[key] = strings.Join(title, ",")
	}
}

func getMdContentWithTitle(doc *schema.Document) string {
	if doc.MetaData == nil {
		return doc.Content
	}
	title := ""
	list := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	for _, v := range list {
		if d, e := doc.MetaData[v].(string); e && len(d) > 0 {
			title += fmt.Sprintf("%s:%s ", v, d)
		}
	}
	if len(title) == 0 {
		return doc.Content
	}
	return title + "\n" + doc.Content
}
