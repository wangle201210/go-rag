package indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/wangle201210/go-rag/server/core/common"
)

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
