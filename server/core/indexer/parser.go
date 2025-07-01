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
				"##": "headerNameOfLevel2",
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

		g.Log().Infof(ctx, "content: %+v", content)

		for _, c := range content {
			doc.Content += c.Content
		}
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

// splitCsvToChunks 解析 csv 内容为单 sheet chunk
func splitCsvToChunks(content string) XlsxSheetChunk {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return XlsxSheetChunk{}
	}
	var headers []string
	var rows [][]string
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := splitCsvLine(line)
		if i == 0 {
			headers = cols
		} else {
			rows = append(rows, cols)
		}
	}
	return XlsxSheetChunk{
		Headers:  headers,
		Rows:     rows,
		StartRow: 2,
		EndRow:   len(rows) + 1,
		ColCount: len(headers),
		Meta:     map[string]interface{}{},
	}
}

// splitCsvLine 简单分割 csv 行（不处理转义）
func splitCsvLine(line string) []string {
	parts := strings.Split(line, ",")
	for i, p := range parts {
		parts[i] = strings.Trim(p, `"`)
	}
	return parts
}

// XlsxParser xlsx 文件解析器（结构化输出）
type XlsxParser struct {
	parser.TextParser
}

// Parse xlsx 解析方法
func (p *XlsxParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	docs, err := p.TextParser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		var parsed struct {
			Sheets []struct {
				Name string     `json:"name"`
				Data [][]string `json:"data"`
			} `json:"sheets"`
		}
		err := json.Unmarshal([]byte(doc.Content), &parsed)
		if err != nil {
			continue
		}
		var chunks []*XlsxSheetChunk
		for _, sheet := range parsed.Sheets {
			if len(sheet.Data) == 0 {
				continue
			}
			headers := sheet.Data[0]
			rows := sheet.Data[1:]
			chunk := &XlsxSheetChunk{
				SheetName: sheet.Name,
				Headers:   headers,
				Rows:      rows,
				StartRow:  2,
				EndRow:    len(rows) + 1,
				ColCount:  len(headers),
				Meta:      map[string]interface{}{},
			}
			chunks = append(chunks, chunk)
		}
		chunkMetas := make([]map[string]interface{}, 0, len(chunks))
		for _, c := range chunks {
			chunkMetas = append(chunkMetas, map[string]interface{}{
				"sheet":    c.SheetName,
				"headers":  c.Headers,
				"rowCount": len(c.Rows),
				"colCount": c.ColCount,
				"startRow": c.StartRow,
				"endRow":   c.EndRow,
				"rows":     c.Rows,
			})
		}
		content, err := json.Marshal(chunkMetas)
		if err != nil {
			return nil, err
		}
		doc.Content = string(content)
	}
	return docs, nil
}

// CSVParser CSV文件解析器（结构化输出，兼容 xlsx）
type CSVParser struct {
	parser.TextParser
}

// Parse CSV解析方法
func (p *CSVParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	docs, err := p.TextParser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		chunks := splitCsvToChunks(doc.Content)
		chunkMetas := map[string]interface{}{
			"headers":  chunks.Headers,
			"rowCount": len(chunks.Rows),
			"colCount": chunks.ColCount,
			"startRow": chunks.StartRow,
			"endRow":   chunks.EndRow,
			"rows":     chunks.Rows,
		}
		content, err := json.Marshal(chunkMetas)
		if err != nil {
			return nil, err
		}
		doc.Content = string(content)
	}
	return docs, nil
}

func newParser(ctx context.Context) (p parser.Parser, err error) {
	// 创建各种解析器实例
	enhancedTextParser := &EnhancedTextParser{}
	markdownParser := &MarkdownParser{}
	csvParser := &CSVParser{}

	// 创建HTML解析器
	htmlParser, err := html.NewParser(ctx, &html.Config{
		Selector: common.Of("body"),
	})
	if err != nil {
		return nil, err
	}

	// 创建PDF解析器
	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		return nil, err
	}

	// 创建解析器
	p, err = parser.NewExtParser(ctx, &parser.ExtParserConfig{
		// 注册特定扩展名的解析器
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
			".md":   markdownParser,
			".csv":  csvParser,
			".txt":  enhancedTextParser,
			".text": enhancedTextParser,
		},
		// 设置默认解析器，用于处理未知格式
		FallbackParser: enhancedTextParser,
	})
	if err != nil {
		return nil, err
	}
	return
}
