package server

var (
	knowledgeBase   IKnowledgeBase
	documentChunk   IDocumentChunk
	documentMapping IDocumentMapping
)

func GetKnowledgeBase() IKnowledgeBase {
	if knowledgeBase == nil {
		knowledgeBase = NewKnowledgeBase()
	}
	return knowledgeBase
}

func GetDocumentChunk() IDocumentChunk {
	if documentChunk == nil {
		documentChunk = NewDocumentChunk()
	}
	return documentChunk
}

func GetDocumentMapping() IDocumentMapping {
	if documentMapping == nil {
		documentMapping = NewDocumentMapping()
	}
	return documentMapping
}
