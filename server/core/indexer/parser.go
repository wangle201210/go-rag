package indexer

import (
	"bytes"
	"context"
	"encoding/csv"
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
	"github.com/xuri/excelize/v2"
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

// XlsxParser xlsx 文件解析器（使用 excelize 库）
type XlsxParser struct {
	parser.Parser
}

// NewXlsxParser 创建 xlsx 解析器
func NewXlsxParser(ctx context.Context) (*XlsxParser, error) {
	p := &parser.TextParser{}
	return &XlsxParser{Parser: p}, nil
}

// Parse 重写 Parse 方法，使用 excelize 库解析所有工作表
func (p *XlsxParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	// 读取整个文件内容到内存
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// 使用 excelize 打开 Excel 文件
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取所有工作表
	sheetList := f.GetSheetList()
	g.Log().Debugf(ctx, "xlsx parser found sheets: %v", sheetList)

	var documents []*schema.Document

	// 遍历所有工作表
	for _, sheetName := range sheetList {
		// 获取当前工作表的所有单元格
		rows, err := f.GetRows(sheetName)
		if err != nil {
			g.Log().Warningf(ctx, "failed to get rows from sheet %s: %v", sheetName, err)
			continue
		}

		// 跳过空工作表
		if len(rows) == 0 {
			continue
		}

		// 提取表头和数据行
		headers := rows[0]
		dataRows := rows[1:]

		// 创建工作表元数据
		sheetData := XlsxSheetChunk{
			SheetName: sheetName,
			Headers:   headers,
			Rows:      dataRows,
			StartRow:  1,
			EndRow:    len(rows),
			ColCount:  len(headers),
			Meta: map[string]interface{}{
				"rowCount": len(dataRows),
				"colCount": len(headers),
			},
		}

		// 将工作表数据转换为 JSON 字符串
		content, err := json.Marshal(sheetData)
		if err != nil {
			g.Log().Warningf(ctx, "failed to marshal sheet data: %v", err)
			continue
		}

		// 创建文档对象
		doc := &schema.Document{
			Content: string(content),
			MetaData: map[string]interface{}{
				"sheetName": sheetName,
				"headers":   headers,
				"rowCount":  len(dataRows),
				"colCount":  len(headers),
			},
		}

		documents = append(documents, doc)
		g.Log().Debugf(ctx, "xlsx parser processed sheet %s with %d rows", sheetName, len(rows))
	}

	if len(documents) == 0 {
		// 如果没有成功解析任何工作表，尝试使用原始解析方法
		return p.Parser.Parse(ctx, bytes.NewReader(data), opts...)
	}

	return documents, nil
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

// Parse 重写 Parse 方法，使用 encoding/csv 包增强 CSV 特定处理
func (p *CSVParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	// 读取整个文件内容到内存
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// 创建 CSV 读取器
	csvReader := csv.NewReader(bytes.NewReader(data))

	// 配置 CSV 读取器
	csvReader.LazyQuotes = true       // 允许字段中的引号
	csvReader.TrimLeadingSpace = true // 修剪前导空格
	csvReader.ReuseRecord = true      // 重用记录以减少内存分配

	// 尝试自动检测分隔符
	if separator := detectSeparator(string(data)); separator != 0 {
		csvReader.Comma = separator
	}

	// 读取所有记录
	records, err := csvReader.ReadAll()
	if err != nil {
		g.Log().Warningf(ctx, "CSV parsing error with standard reader: %v, falling back to original parser", err)
		return p.Parser.Parse(ctx, bytes.NewReader(data), opts...)
	}

	if len(records) == 0 {
		g.Log().Warningf(ctx, "CSV file is empty, falling back to original parser")
		return p.Parser.Parse(ctx, bytes.NewReader(data), opts...)
	}

	// 提取表头和数据行
	headers := records[0]
	dataRows := records[1:]

	// 创建元数据
	csvData := XlsxSheetChunk{
		SheetName: "Sheet1", // CSV 没有工作表概念，默认为 Sheet1
		Headers:   headers,
		Rows:      dataRows,
		StartRow:  1,
		EndRow:    len(records),
		ColCount:  len(headers),
		Meta: map[string]interface{}{
			"rowCount":  len(dataRows),
			"colCount":  len(headers),
			"separator": string(csvReader.Comma),
		},
	}

	// 将数据转换为 JSON 字符串
	content, err := json.Marshal(csvData)
	if err != nil {
		g.Log().Warningf(ctx, "Failed to marshal CSV data: %v", err)
		return p.Parser.Parse(ctx, bytes.NewReader(data), opts...)
	}

	// 创建文档对象
	doc := &schema.Document{
		Content: string(content),
		MetaData: map[string]interface{}{
			"headers":   headers,
			"rowCount":  len(dataRows),
			"colCount":  len(headers),
			"separator": string(csvReader.Comma),
		},
	}

	g.Log().Infof(ctx, "CSV parser processed %d rows with separator '%c'", len(records), csvReader.Comma)

	return []*schema.Document{doc}, nil
}

// detectSeparator 尝试检测 CSV 文件的分隔符
func detectSeparator(content string) rune {
	// 常见的分隔符
	separators := []rune{',', ';', '\t', '|'}

	// 获取第一行
	firstLine := ""
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine = lines[0]
	} else {
		return 0
	}

	// 计算每个分隔符在第一行中出现的次数
	counts := make(map[rune]int)
	for _, sep := range separators {
		counts[sep] = strings.Count(firstLine, string(sep))
	}

	// 找到出现次数最多的分隔符
	maxCount := 0
	var bestSeparator rune
	for sep, count := range counts {
		if count > maxCount {
			maxCount = count
			bestSeparator = sep
		}
	}

	// 如果找到合适的分隔符，返回它
	if maxCount > 0 {
		return bestSeparator
	}

	// 默认使用逗号
	return ','
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
