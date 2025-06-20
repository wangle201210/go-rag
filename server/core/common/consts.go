package common

const (
	FieldContent         = "content"
	FieldContentVector   = "content_vector"
	FieldQAContent       = "qa_content"
	FieldQAContentVector = "qa_content_vector"
	FieldExtra           = "ext"
	KnowledgeName        = "_knowledge_name"

	RetrieverFieldKey = "_retriever_field"

	Title1 = "h1"
	Title2 = "h2"
	Title3 = "h3"
)

var (
	// ExtKeys ext 里面需要存储的数据
	ExtKeys = []string{"_extension", "_file_name", "_source", Title1, Title2, Title3}
)
