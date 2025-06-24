// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KnowledgeDocuments is the golang structure for table knowledge_documents.
type KnowledgeDocuments struct {
	Id                int64       `json:"id"                orm:"id"                  description:""` //
	KnowledgeBaseName string      `json:"knowledgeBaseName" orm:"knowledge_base_name" description:""` //
	FileName          string      `json:"fileName"          orm:"file_name"           description:""` //
	Status            int         `json:"status"            orm:"status"              description:""` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""` //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""` //
}
