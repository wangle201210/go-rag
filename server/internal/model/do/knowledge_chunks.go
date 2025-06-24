// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KnowledgeChunks is the golang structure of table knowledge_chunks for DAO operations like Where/Data.
type KnowledgeChunks struct {
	g.Meta         `orm:"table:knowledge_chunks, do:true"`
	Id             interface{} //
	KnowledgeDocId interface{} //
	ChunkId        interface{} //
	Content        interface{} //
	Ext            interface{} //
	Status         interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
