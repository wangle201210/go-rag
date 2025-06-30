package indexer

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/document/parser"
	"github.com/stretchr/testify/assert"
)

func TestEnhancedTextParser(t *testing.T) {
	parser := &EnhancedTextParser{}

	// 测试文本清理
	content := "  这是测试文本  \n\n\n  包含多余的空行  \n\n  和空格  "
	cleaned := parser.cleanAndNormalizeText(content)

	// 验证清理结果
	assert.NotContains(t, cleaned, "   ")
	assert.NotContains(t, cleaned, "\n\n\n")

	// 测试文档质量评估
	score := parser.assessDocumentQuality("这是一个包含中文的测试文档。")
	assert.Greater(t, score, 0.0)
	assert.LessOrEqual(t, score, 1.0)

	// 测试词数统计
	wordCount := parser.countWords("这是 一个 测试 文档")
	assert.Equal(t, 4, wordCount)

	// 测试语言检测
	language := parser.detectLanguage("这是中文文档")
	assert.Equal(t, "zh", language)

	language = parser.detectLanguage("This is English document")
	assert.Equal(t, "en", language)
}

func TestMarkdownParser(t *testing.T) {
	parser := &MarkdownParser{}

	markdownContent := `# 标题1

## 标题2

### 标题3

这是正文内容。

- 列表项1
- 列表项2

代码块

[链接文本](http://example.com)
`

	// 测试结构提取
	structure := parser.extractMarkdownStructure(markdownContent)
	assert.NotNil(t, structure)

	sections, ok := structure["sections"].([]map[string]string)
	assert.True(t, ok)
	assert.Greater(t, len(sections), 0)

	// 测试内容清理
	cleaned := parser.cleanMarkdownContent(markdownContent)
	assert.NotContains(t, cleaned, "#")
	assert.NotContains(t, cleaned, "- ")
	assert.Contains(t, cleaned, "标题1")
	assert.Contains(t, cleaned, "这是正文内容")
}

func TestCSVParser(t *testing.T) {
	parser := &CSVParser{}

	csvContent := `姓名,年龄,城市
张三,25,北京
李四,30,上海
王五,28,广州`

	// 测试CSV结构解析
	structure := parser.parseCSVStructure(csvContent)
	assert.NotNil(t, structure)

	headers, ok := structure["headers"].([]string)
	assert.True(t, ok)
	assert.Equal(t, 3, len(headers))
	assert.Equal(t, "姓名", headers[0])

	rowCount, ok := structure["row_count"].(int)
	assert.True(t, ok)
	assert.Equal(t, 3, rowCount)

	columns, ok := structure["columns"].(int)
	assert.True(t, ok)
	assert.Equal(t, 3, columns)
}

func TestJSONParser(t *testing.T) {
	parser := &JSONParser{}

	// 测试对象类型JSON
	objectJSON := `{"name": "张三", "age": 25, "city": "北京"}`
	structure := parser.parseJSONStructure(objectJSON)
	assert.NotNil(t, structure)

	jsonType, ok := structure["type"].(string)
	assert.True(t, ok)
	assert.Equal(t, "object", jsonType)

	// 测试数组类型JSON
	arrayJSON := `[{"name": "张三"}, {"name": "李四"}]`
	structure = parser.parseJSONStructure(arrayJSON)
	assert.NotNil(t, structure)

	jsonType, ok = structure["type"].(string)
	assert.True(t, ok)
	assert.Equal(t, "array", jsonType)
}

func TestSmartParser(t *testing.T) {
	// 创建智能解析器
	enhancedTextParser := &EnhancedTextParser{}
	markdownParser := &MarkdownParser{}
	csvParser := &CSVParser{}
	jsonParser := &JSONParser{}

	_ = &SmartParser{
		parsers: map[string]parser.Parser{
			".md":   markdownParser,
			".csv":  csvParser,
			".json": jsonParser,
			".txt":  enhancedTextParser,
		},
		fallbackParser: enhancedTextParser,
	}

	// 测试文件扩展名获取
	ext := getFileExtension("test.md")
	assert.Equal(t, ".md", ext)

	ext = getFileExtension("test.CSV")
	assert.Equal(t, ".csv", ext)

	// 测试文件类型判断
	assert.True(t, isMarkdownFile("test.md"))
	assert.True(t, isMarkdownFile("test.markdown"))
	assert.False(t, isMarkdownFile("test.txt"))

	assert.True(t, isCSVFile("test.csv"))
	assert.False(t, isCSVFile("test.txt"))

	assert.True(t, isJSONFile("test.json"))
	assert.False(t, isJSONFile("test.txt"))
}

func TestParserIntegration(t *testing.T) {
	ctx := context.Background()

	// 测试解析器创建
	parser, err := newParser(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, parser)

	// 验证解析器类型
	_, ok := parser.(*SmartParser)
	assert.True(t, ok, "解析器应该是SmartParser类型")
}
