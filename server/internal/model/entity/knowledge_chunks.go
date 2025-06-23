// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KnowledgeChunks is the golang structure for table knowledge_chunks.
type KnowledgeChunks struct {
	Id             int         `json:"id"             orm:"id"               description:""` //
	KnowledgeDocId int         `json:"knowledgeDocId" orm:"knowledge_doc_id" description:""` //
	EsChunkId      string      `json:"esChunkId"      orm:"es_chunk_id"      description:""` //
	Content        string      `json:"content"        orm:"content"          description:""` //
	Ext            string      `json:"ext"            orm:"ext"              description:""` //
	Status         int         `json:"status"         orm:"status"           description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:""` //
}
