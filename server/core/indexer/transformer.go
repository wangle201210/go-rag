package indexer

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

// newDocumentTransformer component initialization function of node 'DocumentTransformer3' in graph 'rag'
func newDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	trans := &transformer{}
	// 递归分割
	config := &recursive.Config{
		ChunkSize:   1000, // 每段内容1000字
		OverlapSize: 100,  // 有10%的重叠
		Separators:  []string{"\n", "。", "?", "？", "!", "！"},
	}
	recTrans, err := recursive.NewSplitter(ctx, config)
	if err != nil {
		return nil, err
	}
	// md 文档特殊处理
	mdTrans, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers:     map[string]string{"#": "h1", "##": "h2", "###": "h3"},
		TrimHeaders: false,
	})
	if err != nil {
		return nil, err
	}
	trans.recursive = recTrans
	trans.markdown = mdTrans
	return trans, nil
}

type transformer struct {
	markdown  document.Transformer
	recursive document.Transformer
}

func (x *transformer) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	isMd := false
	for _, doc := range docs {
		// 只需要判断第一个是不是.md
		if doc.MetaData["_extension"] == ".md" {
			isMd = true
			break
		}
	}
	if isMd {
		return x.markdown.Transform(ctx, docs, opts...)
	}
	return x.recursive.Transform(ctx, docs, opts...)
}

//	type Chunk struct {
//		ID       string
//		Content  string
//		Metadata map[string]interface{}
//	}
//
//	type Transformer struct {
//		chunkSize int
//	}
//
//	func NewTransformer(chunkSize int) *Transformer {
//		return &Transformer{
//			chunkSize: chunkSize,
//		}
//	}
//
//	func (t *Transformer) Transform(doc *Document) ([]*Chunk, error) {
//		// 将文档内容分割成块
//		chunks := t.splitIntoChunks(doc.Content)
//
//		// 为每个块创建Chunk对象
//		var result []*Chunk
//		for i, content := range chunks {
//			chunk := &Chunk{
//				ID:      generateChunkID(doc.ID, i),
//				Content: content,
//				Metadata: map[string]interface{}{
//					"document_id":   doc.ID,
//					"document_name": doc.Name,
//					"document_type": doc.Type,
//					"chunk_index":   i,
//				},
//			}
//			result = append(result, chunk)
//		}
//
//		return result, nil
//	}
//
//	func (t *Transformer) splitIntoChunks(content string) []string {
//		// 按段落分割内容
//		paragraphs := strings.Split(content, "\n\n")
//
//		var chunks []string
//		var currentChunk strings.Builder
//
//		for _, paragraph := range paragraphs {
//			// 如果当前块加上新段落超过chunkSize，先保存当前块
//			if currentChunk.Len()+len(paragraph) > t.chunkSize {
//				if currentChunk.Len() > 0 {
//					chunks = append(chunks, currentChunk.String())
//					currentChunk.Reset()
//				}
//			}
//
//			// 如果单个段落超过chunkSize，需要进一步分割
//			if len(paragraph) > t.chunkSize {
//				sentences := strings.Split(paragraph, ". ")
//				for _, sentence := range sentences {
//					if currentChunk.Len()+len(sentence) > t.chunkSize {
//						if currentChunk.Len() > 0 {
//							chunks = append(chunks, currentChunk.String())
//							currentChunk.Reset()
//						}
//					}
//					currentChunk.WriteString(sentence)
//					currentChunk.WriteString(". ")
//				}
//			} else {
//				currentChunk.WriteString(paragraph)
//				currentChunk.WriteString("\n\n")
//			}
//		}
//
//		// 添加最后一个块
//		if currentChunk.Len() > 0 {
//			chunks = append(chunks, currentChunk.String())
//		}
//
//		return chunks
//	}
func generateChunkID(docID string, index int) string {
	return fmt.Sprintf("%s_chunk_%d", docID, index)
}
