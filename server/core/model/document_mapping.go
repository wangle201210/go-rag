package model

type DocumentMapping struct {
	Id              int64  `json:"id"`
	KnowledgeBaseId int64  `json:"knowledgeBaseId"`
	DocumentId      string `json:"documentId"`
	DocumentName    string `json:"documentName"`
	DocumentType    string `json:"documentType"`
	DocumentPath    string `json:"documentPath"`
	DocumentSize    int64  `json:"documentSize"`
	ChunkCount      int    `json:"chunkCount"`
	Status          int    `json:"status"`
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
}

type DocumentMappingCreateReq struct {
	KnowledgeBaseId string `json:"knowledgeBaseId" v:"required#知识库ID不能为空"`
	DocumentId      string `json:"documentId" v:"required#文档ID不能为空"`
	DocumentName    string `json:"documentName" v:"required#文档名称不能为空"`
	DocumentType    string `json:"documentType" v:"required|in:md,pdf,html#文档类型不能为空|文档类型不正确"`
	DocumentPath    string `json:"documentPath" v:"required#文档路径不能为空"`
	DocumentSize    int64  `json:"documentSize" v:"required#文档大小不能为空"`
	ChunkCount      int    `json:"chunkCount" v:"required#分块数量不能为空"`
}

type DocumentMappingCreateRes struct {
	Id int64 `json:"id"`
}

type DocumentMappingUpdateReq struct {
	Id              int64  `json:"id" v:"required#映射ID不能为空"`
	KnowledgeBaseId int64  `json:"knowledgeBaseId" v:"required#知识库ID不能为空"`
	DocumentName    string `json:"documentName" v:"required#文档名称不能为空"`
	DocumentType    string `json:"documentType" v:"required|in:md,pdf,html#文档类型不能为空|文档类型不正确"`
	DocumentPath    string `json:"documentPath" v:"required#文档路径不能为空"`
	DocumentSize    int64  `json:"documentSize" v:"required#文档大小不能为空"`
	ChunkCount      int    `json:"chunkCount" v:"required#分块数量不能为空"`
	Status          int    `json:"status" v:"required|in:0,1#状态不能为空|状态值不正确"`
}

type DocumentMappingUpdateRes struct {
	Success bool `json:"success"`
}

type DocumentMappingDeleteReq struct {
	Id int64 `json:"id" v:"required#映射ID不能为空"`
}

type DocumentMappingDeleteRes struct {
	Success bool `json:"success"`
}

type DocumentMappingListReq struct {
	Page            int   `json:"page" v:"required|min:1#页码不能为空|页码最小值为1"`
	PageSize        int   `json:"pageSize" v:"required|min:1#每页数量不能为空|每页数量最小值为1"`
	KnowledgeBaseId int64 `json:"knowledgeBaseId,optional"`
	Status          int   `json:"status,optional"`
}

type DocumentMappingListRes struct {
	List  []DocumentMapping `json:"list"`
	Total int               `json:"total"`
}
