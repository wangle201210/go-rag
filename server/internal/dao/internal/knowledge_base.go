// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KnowledgeBaseDao is the data access object for the table knowledge_base.
type KnowledgeBaseDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of the current DAO.
	columns KnowledgeBaseColumns // columns contains all the column names of Table for convenient usage.
}

// KnowledgeBaseColumns defines and stores column names for the table knowledge_base.
type KnowledgeBaseColumns struct {
	Id          string // 主键ID
	Name        string // 知识库名称
	Description string // 知识库描述
	Category    string // 知识库分类
	Status      string // 状态：1-启用,2-禁用
	CreateTime  string // 创建时间
	UpdateTime  string // 更新时间
}

// knowledgeBaseColumns holds the columns for the table knowledge_base.
var knowledgeBaseColumns = KnowledgeBaseColumns{
	Id:          "id",
	Name:        "name",
	Description: "description",
	Category:    "category",
	Status:      "status",
	CreateTime:  "create_time",
	UpdateTime:  "update_time",
}

// NewKnowledgeBaseDao creates and returns a new DAO object for table data access.
func NewKnowledgeBaseDao() *KnowledgeBaseDao {
	return &KnowledgeBaseDao{
		group:   "default",
		table:   "knowledge_base",
		columns: knowledgeBaseColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KnowledgeBaseDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KnowledgeBaseDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KnowledgeBaseDao) Columns() KnowledgeBaseColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KnowledgeBaseDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KnowledgeBaseDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *KnowledgeBaseDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
