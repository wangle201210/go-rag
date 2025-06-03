package indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/wangle201210/go-rag/server/core/common"
)

//
// type Document struct {
// 	ID      string
// 	Name    string
// 	Type    string
// 	Content string
// }
//
// type Parser struct{}
//
// func NewParser() *Parser {
// 	return &Parser{}
// }
//
// func (p *Parser) Parse(content string) (*Document, error) {
// 	// 解析文档内容
// 	// 这里可以根据不同的文档类型实现不同的解析逻辑
// 	// 目前简单实现，后续可以扩展
//
// 	// 生成文档ID
// 	docID := generateDocID(content)
//
// 	// 获取文档名称
// 	docName := extractDocName(content)
//
// 	// 获取文档类型
// 	docType := determineDocType(content)
//
// 	return &Document{
// 		ID:      docID,
// 		Name:    docName,
// 		Type:    docType,
// 		Content: content,
// 	}, nil
// }
//
// func generateDocID(content string) string {
// 	// 使用内容的前32个字符作为文档ID
// 	// 实际应用中应该使用更复杂的算法，如MD5等
// 	if len(content) > 32 {
// 		return content[:32]
// 	}
// 	return content
// }
//
// func extractDocName(content string) string {
// 	// 从内容中提取文档名称
// 	// 这里简单实现，取第一行作为文档名称
// 	lines := strings.Split(content, "\n")
// 	if len(lines) > 0 {
// 		return strings.TrimSpace(lines[0])
// 	}
// 	return "未命名文档"
// }
//
// func determineDocType(content string) string {
// 	// 根据内容特征判断文档类型
// 	// 这里简单实现，实际应用中应该更复杂
// 	if strings.Contains(content, "# ") {
// 		return "md"
// 	}
// 	if strings.Contains(content, "<html") {
// 		return "html"
// 	}
// 	return "txt"
// }

func newParser(ctx context.Context) (p parser.Parser, err error) {
	textParser := parser.TextParser{}

	htmlParser, err := html.NewParser(ctx, &html.Config{
		Selector: common.Of("body"),
	})
	if err != nil {
		return nil, err
	}

	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		return
	}

	// 创建扩展解析器
	p, err = parser.NewExtParser(ctx, &parser.ExtParserConfig{
		// 注册特定扩展名的解析器
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
		},
		// 设置默认解析器，用于处理未知格式
		FallbackParser: textParser,
	})
	if err != nil {
		return nil, err
	}
	return
}
