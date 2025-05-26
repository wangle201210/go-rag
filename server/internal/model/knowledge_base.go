package model

type KnowledgeBaseListReq struct {
	Page     int    `json:"page" v:"required|min:1#页码不能为空|页码最小值为1"`
	PageSize int    `json:"pageSize" v:"required|min:1#每页数量不能为空|每页数量最小值为1"`
	Category string `json:"category,optional"`
	Keyword  string `json:"keyword,optional"`
}

type KnowledgeBaseListRes struct {
	List  []KnowledgeBaseItem `json:"list"`
	Total int                 `json:"total"`
}

type KnowledgeBaseItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Status      int    `json:"status"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

type KnowledgeBaseCreateReq struct {
	Name        string `json:"name" v:"required|length:1,50#知识库名称不能为空|知识库名称长度应在1-50之间"`
	Description string `json:"description" v:"required|length:1,200#知识库描述不能为空|知识库描述长度应在1-200之间"`
	Category    string `json:"category" v:"required|length:1,50#知识库分类不能为空|知识库分类长度应在1-50之间"`
}

type KnowledgeBaseCreateRes struct {
	Id int64 `json:"id"`
}

type KnowledgeBaseUpdateReq struct {
	Id          int64  `json:"id" v:"required#知识库ID不能为空"`
	Name        string `json:"name" v:"required|length:1,50#知识库名称不能为空|知识库名称长度应在1-50之间"`
	Description string `json:"description" v:"required|length:1,200#知识库描述不能为空|知识库描述长度应在1-200之间"`
	Category    string `json:"category" v:"required|length:1,50#知识库分类不能为空|知识库分类长度应在1-50之间"`
	Status      int    `json:"status" v:"required|in:0,1#状态不能为空|状态值不正确"`
}

type KnowledgeBaseUpdateRes struct {
	Success bool `json:"success"`
}

type KnowledgeBaseDeleteReq struct {
	Id int64 `json:"id" v:"required#知识库ID不能为空"`
}

type KnowledgeBaseDeleteRes struct {
	Success bool `json:"success"`
}
