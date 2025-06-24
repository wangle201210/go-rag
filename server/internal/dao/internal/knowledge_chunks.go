// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KnowledgeChunksDao is the data access object for the table knowledge_chunks.
type KnowledgeChunksDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  KnowledgeChunksColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// KnowledgeChunksColumns defines and stores column names for the table knowledge_chunks.
type KnowledgeChunksColumns struct {
	Id             string //
	KnowledgeDocId string //
	ChunkId        string //
	Content        string //
	Ext            string //
	Status         string //
	CreatedAt      string //
	UpdatedAt      string //
}

// knowledgeChunksColumns holds the columns for the table knowledge_chunks.
var knowledgeChunksColumns = KnowledgeChunksColumns{
	Id:             "id",
	KnowledgeDocId: "knowledge_doc_id",
	ChunkId:        "chunk_id",
	Content:        "content",
	Ext:            "ext",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewKnowledgeChunksDao creates and returns a new DAO object for table data access.
func NewKnowledgeChunksDao(handlers ...gdb.ModelHandler) *KnowledgeChunksDao {
	return &KnowledgeChunksDao{
		group:    "default",
		table:    "knowledge_chunks",
		columns:  knowledgeChunksColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KnowledgeChunksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KnowledgeChunksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KnowledgeChunksDao) Columns() KnowledgeChunksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KnowledgeChunksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KnowledgeChunksDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *KnowledgeChunksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
