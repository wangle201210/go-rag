// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KnowledgeDocumentsDao is the data access object for the table knowledge_documents.
type KnowledgeDocumentsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KnowledgeDocumentsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KnowledgeDocumentsColumns defines and stores column names for the table knowledge_documents.
type KnowledgeDocumentsColumns struct {
	Id                string //
	KnowledgeBaseName string //
	FileName          string //
	Status            string //
	CreatedAt         string //
	UpdatedAt         string //
}

// knowledgeDocumentsColumns holds the columns for the table knowledge_documents.
var knowledgeDocumentsColumns = KnowledgeDocumentsColumns{
	Id:                "id",
	KnowledgeBaseName: "knowledge_base_name",
	FileName:          "file_name",
	Status:            "status",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewKnowledgeDocumentsDao creates and returns a new DAO object for table data access.
func NewKnowledgeDocumentsDao(handlers ...gdb.ModelHandler) *KnowledgeDocumentsDao {
	return &KnowledgeDocumentsDao{
		group:    "default",
		table:    "knowledge_documents",
		columns:  knowledgeDocumentsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KnowledgeDocumentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KnowledgeDocumentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KnowledgeDocumentsDao) Columns() KnowledgeDocumentsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KnowledgeDocumentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KnowledgeDocumentsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KnowledgeDocumentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
