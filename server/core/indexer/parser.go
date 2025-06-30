package indexer

import (
	"context"
	"io"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
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

		// 添加文档质量评估
		doc.MetaData["quality_score"] = p.assessDocumentQuality(doc.Content)
		doc.MetaData["word_count"] = p.countWords(doc.Content)
		doc.MetaData["language"] = p.detectLanguage(doc.Content)
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

// assessDocumentQuality 评估文档质量
func (p *EnhancedTextParser) assessDocumentQuality(text string) float64 {
	if len(text) == 0 {
		return 0.0
	}

	score := 1.0

	// 检查文本长度
	if len(text) < 50 {
		score *= 0.5
	} else if len(text) > 10000 {
		score *= 0.9
	}

	// 检查中文字符比例
	chineseCount := 0
	totalChars := 0
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			chineseCount++
		}
		if !unicode.IsSpace(r) {
			totalChars++
		}
	}

	if totalChars > 0 {
		chineseRatio := float64(chineseCount) / float64(totalChars)
		if chineseRatio > 0.1 { // 包含中文
			score *= 1.1
		}
	}

	// 检查标点符号使用
	sentenceCount := len(strings.Split(text, "。")) + len(strings.Split(text, "！")) + len(strings.Split(text, "？"))
	if sentenceCount > 0 {
		avgSentenceLength := float64(len(text)) / float64(sentenceCount)
		if avgSentenceLength > 200 {
			score *= 0.8 // 句子过长
		}
	}

	return score
}

// countWords 统计词数
func (p *EnhancedTextParser) countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

// detectLanguage 检测语言
func (p *EnhancedTextParser) detectLanguage(text string) string {
	chineseCount := 0
	englishCount := 0

	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			chineseCount++
		} else if unicode.IsLetter(r) {
			englishCount++
		}
	}

	if chineseCount > englishCount {
		return "zh"
	} else if englishCount > chineseCount {
		return "en"
	}
	return "mixed"
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
		// 提取Markdown标题结构
		doc.MetaData["markdown_structure"] = p.extractMarkdownStructure(doc.Content)
		// 清理Markdown标记
		doc.Content = p.cleanMarkdownContent(doc.Content)
	}

	return docs, nil
}

// extractMarkdownStructure 提取Markdown结构
func (p *MarkdownParser) extractMarkdownStructure(content string) map[string]interface{} {
	structure := make(map[string]interface{})
	lines := strings.Split(content, "\n")

	var currentH1, currentH2, currentH3 string
	var sections []map[string]string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			currentH1 = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			currentH2 = ""
			currentH3 = ""
		} else if strings.HasPrefix(line, "## ") {
			currentH2 = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			currentH3 = ""
		} else if strings.HasPrefix(line, "### ") {
			currentH3 = strings.TrimSpace(strings.TrimPrefix(line, "### "))
		}

		if currentH1 != "" {
			sections = append(sections, map[string]string{
				"h1": currentH1,
				"h2": currentH2,
				"h3": currentH3,
			})
		}
	}

	structure["sections"] = sections
	return structure
}

// cleanMarkdownContent 清理Markdown内容
func (p *MarkdownParser) cleanMarkdownContent(content string) string {
	// 移除Markdown标记但保留文本内容
	lines := strings.Split(content, "\n")
	var cleanedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 移除标题标记
		if strings.HasPrefix(line, "#") {
			line = strings.TrimSpace(strings.TrimLeft(line, "# "))
		}

		// 移除列表标记
		if strings.HasPrefix(line, "- ") {
			line = strings.TrimPrefix(line, "- ")
		}
		if strings.HasPrefix(line, "* ") {
			line = strings.TrimPrefix(line, "* ")
		}
		if strings.HasPrefix(line, "+ ") {
			line = strings.TrimPrefix(line, "+ ")
		}

		// 移除代码块标记
		if strings.HasPrefix(line, "```") {
			continue
		}

		// 移除行内代码标记
		line = strings.ReplaceAll(line, "`", "")

		// 移除链接标记，保留文本
		if strings.Contains(line, "[") && strings.Contains(line, "](") {
			// 简单的链接处理
			start := strings.Index(line, "[")
			end := strings.Index(line, "]")
			if start != -1 && end != -1 && end > start {
				linkText := line[start+1 : end]
				line = strings.Replace(line, line[start:strings.Index(line, ")")+1], linkText, 1)
			}
		}

		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// CSVParser CSV文件解析器
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
		// 解析CSV结构
		csvData := p.parseCSVStructure(doc.Content)
		doc.MetaData["csv_structure"] = csvData
		doc.MetaData["file_type"] = "csv"
	}

	return docs, nil
}

// parseCSVStructure 解析CSV结构
func (p *CSVParser) parseCSVStructure(content string) map[string]interface{} {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return map[string]interface{}{}
	}

	// 提取表头
	headers := strings.Split(lines[0], ",")
	for i, header := range headers {
		headers[i] = strings.TrimSpace(strings.Trim(header, `"`))
	}

	// 统计行数
	rowCount := len(lines) - 1 // 减去表头

	return map[string]interface{}{
		"headers":   headers,
		"row_count": rowCount,
		"columns":   len(headers),
	}
}

// JSONParser JSON文件解析器
type JSONParser struct {
	parser.TextParser
}

// Parse JSON解析方法
func (p *JSONParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	docs, err := p.TextParser.Parse(ctx, reader, opts...)
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		// 解析JSON结构
		jsonStructure := p.parseJSONStructure(doc.Content)
		doc.MetaData["json_structure"] = jsonStructure
		doc.MetaData["file_type"] = "json"
	}

	return docs, nil
}

// parseJSONStructure 解析JSON结构
func (p *JSONParser) parseJSONStructure(content string) map[string]interface{} {
	// 简单的JSON结构分析
	structure := make(map[string]interface{})

	// 检测是否为数组
	if strings.TrimSpace(content)[0] == '[' {
		structure["type"] = "array"
		// 计算数组元素数量（简单估算）
		structure["element_count"] = strings.Count(content, "{")
	} else {
		structure["type"] = "object"
		// 计算顶级字段数量
		structure["field_count"] = strings.Count(content, `":`)
	}

	return structure
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// isMarkdownFile 判断是否为Markdown文件
func isMarkdownFile(filename string) bool {
	ext := getFileExtension(filename)
	return ext == ".md" || ext == ".markdown"
}

// isCSVFile 判断是否为CSV文件
func isCSVFile(filename string) bool {
	ext := getFileExtension(filename)
	return ext == ".csv"
}

// isJSONFile 判断是否为JSON文件
func isJSONFile(filename string) bool {
	ext := getFileExtension(filename)
	return ext == ".json"
}

func newParser(ctx context.Context) (p parser.Parser, err error) {
	// 创建各种解析器实例
	enhancedTextParser := &EnhancedTextParser{}
	markdownParser := &MarkdownParser{}
	csvParser := &CSVParser{}
	jsonParser := &JSONParser{}

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

	// 创建智能解析器
	smartParser := &SmartParser{
		parsers: map[string]parser.Parser{
			".html":     htmlParser,
			".htm":      htmlParser,
			".pdf":      pdfParser,
			".md":       markdownParser,
			".markdown": markdownParser,
			".csv":      csvParser,
			".json":     jsonParser,
			".txt":      enhancedTextParser,
			".text":     enhancedTextParser,
		},
		fallbackParser: enhancedTextParser,
	}

	return smartParser, nil
}

// SmartParser 智能解析器，根据文件扩展名选择合适的解析器
type SmartParser struct {
	parsers        map[string]parser.Parser
	fallbackParser parser.Parser
}

// Parse 智能解析方法
func (sp *SmartParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	// 由于无法从reader获取文件名，这里使用默认解析器
	// 在实际使用中，文件名应该通过其他方式传递
	g.Log().Infof(ctx, "使用默认解析器处理文件")
	return sp.fallbackParser.Parse(ctx, reader, opts...)
}
