package indexer

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/common"
)

// EnhancedTextParser 增强的文本解析器
type EnhancedTextParser struct {
	parser.TextParser
}

// Parse 增强的文本解析方法
func (p *EnhancedTextParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	// 先调用原始解析器
	docs, err := p.TextParser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}

	// 对解析结果进行后处理
	for _, doc := range docs {
		// 清理和标准化文本内容
		doc.Content = p.cleanAndNormalizeText(doc.Content)
		g.Log().Debugf(ctx, "enhanced text parser doc Content: %+v", doc.Content)
	}

	return docs, nil
}

// cleanAndNormalizeText 清理和标准化文本
func (p *EnhancedTextParser) cleanAndNormalizeText(text string) string {
	// 移除多余的空白字符
	text = strings.TrimSpace(text)

	// 标准化换行符
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 移除连续的空行
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	prevEmpty := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if !prevEmpty {
				cleanedLines = append(cleanedLines, line)
			}
			prevEmpty = true
		} else {
			cleanedLines = append(cleanedLines, line)
			prevEmpty = false
		}
	}

	// 移除开头和结尾的空行
	if len(cleanedLines) > 0 && cleanedLines[0] == "" {
		cleanedLines = cleanedLines[1:]
	}
	if len(cleanedLines) > 0 && cleanedLines[len(cleanedLines)-1] == "" {
		cleanedLines = cleanedLines[:len(cleanedLines)-1]
	}

	return strings.Join(cleanedLines, "\n")
}

// MarkdownParser Markdown专用解析器
type MarkdownParser struct {
	parser.TextParser
}

// Parse Markdown解析方法
func (p *MarkdownParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	docs, err := p.TextParser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
			Headers: map[string]string{
				"##":  "headerNameOfLevel2",
				"###": "headerNameOfLevel3",
			},
		})
		if err != nil {
			return nil, err
		}

		content, err := splitter.Transform(ctx, []*schema.Document{
			{
				Content: doc.Content,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, c := range content {
			doc.Content += c.Content
		}
		g.Log().Debugf(ctx, "markdown parser doc Content: %+v", doc.Content)
	}

	return docs, nil
}

// XlsxSheetChunk 表示一个 xlsx/csv 的 sheet 分块
type XlsxSheetChunk struct {
	SheetName string
	Headers   []string
	Rows      [][]string
	StartRow  int
	EndRow    int
	ColCount  int
	Meta      map[string]interface{}
}

// XlsxParser xlsx 文件解析器（使用 eino-ext 官方库）
type XlsxParser struct {
	parser.Parser
}

// NewXlsxParser 创建 xlsx 解析器
func NewXlsxParser(ctx context.Context) (*XlsxParser, error) {
	p := &parser.TextParser{}
	return &XlsxParser{Parser: p}, nil
}

// Parse 重写 Parse 方法，增加后处理
func (p *XlsxParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	docs, err := p.Parser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var rows [][]string
		if err := json.Unmarshal([]byte(doc.Content), &rows); err != nil {
			var parsed struct {
				Sheets []struct {
					Name string     `json:"name"`
					Data [][]string `json:"data"`
				} `json:"sheets"`
			}
			if err := json.Unmarshal([]byte(doc.Content), &parsed); err == nil {
				for _, sheet := range parsed.Sheets {
					if len(sheet.Data) > 0 {
						rows = sheet.Data
						break
					}
				}
			}
		}

		if len(rows) > 0 {
			headers := rows[0]
			dataRows := rows[1:]

			chunkMeta := map[string]interface{}{
				"headers":  headers,
				"rowCount": len(dataRows),
				"colCount": len(headers),
				"rows":     dataRows,
			}

			content, err := json.Marshal(chunkMeta)
			if err != nil {
				return nil, err
			}
			doc.Content = string(content)
		}
		g.Log().Debugf(ctx, "xlsx parser doc Content: %+v", doc.Content)
	}

	return docs, nil
}

// CSVParser CSV文件解析器（结构化处理）
type CSVParser struct {
	parser.Parser
}

// NewCSVParser 创建 CSV 解析器
func NewCSVParser(ctx context.Context) (*CSVParser, error) {
	p := &parser.TextParser{}
	return &CSVParser{Parser: p}, nil
}

// Parse 重写 Parse 方法，增加 CSV 特定处理
func (p *CSVParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	docs, err := p.Parser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		rows := parseCSVContent(doc.Content)

		if len(rows) > 0 {
			headers := rows[0]
			dataRows := rows[1:]

			chunkMeta := map[string]interface{}{
				"headers":  headers,
				"rowCount": len(dataRows),
				"colCount": len(headers),
				"rows":     dataRows,
			}

			content, err := json.Marshal(chunkMeta)
			if err != nil {
				return nil, err
			}
			doc.Content = string(content)
		}
		g.Log().Infof(ctx, "csv parser doc Content: %+v", doc.Content)
	}

	return docs, nil
}

// parseCSVContent 解析 CSV 内容为行列结构
func parseCSVContent(content string) [][]string {
	lines := strings.Split(content, "\n")
	var rows [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row []string
		for _, cell := range strings.Split(line, ",") {
			row = append(row, strings.Trim(cell, `"`))
		}
		rows = append(rows, row)
	}
	return rows
}

func newParser(ctx context.Context) (p parser.Parser, err error) {
	enhancedTextParser := &EnhancedTextParser{}
	markdownParser := &MarkdownParser{}

	xlsxParser, err := NewXlsxParser(ctx)
	if err != nil {
		return nil, err
	}

	csvParser, err := NewCSVParser(ctx)
	if err != nil {
		return nil, err
	}

	htmlParser, err := html.NewParser(ctx, &html.Config{
		Selector: common.Of("body"),
	})
	if err != nil {
		return nil, err
	}

	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		return nil, err
	}

	p, err = parser.NewExtParser(ctx, &parser.ExtParserConfig{
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
			".md":   markdownParser,
			".csv":  csvParser,
			".xlsx": xlsxParser,
			".xls":  xlsxParser,
			".txt":  enhancedTextParser,
			".text": enhancedTextParser,
		},
		FallbackParser: enhancedTextParser,
	})
	if err != nil {
		return nil, err
	}
	return
}
