package knowledge

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/internal/dao"
	"github.com/wangle201210/go-rag/server/internal/model/entity"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// SaveDocumentsInfo 保存文档信息
func SaveDocumentsInfo(ctx context.Context, documents entity.KnowledgeDocuments) (id int64, err error) {
	// 确保 ID 为 0，让数据库自动分配
	documents.Id = 0

	// OmitEmpty 会忽略零值字段（包括 ID=0），让数据库自动分配 ID
	// 这样可以兼容 MySQL 和 SQLite 的自增主键
	result, err := dao.KnowledgeDocuments.Ctx(ctx).Data(documents).OmitEmpty().Insert()
	if err != nil {
		g.Log().Errorf(ctx, "保存文档信息失败: %+v, 错误: %v", documents, err)
		return 0, fmt.Errorf("保存文档信息失败: %w", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}
	g.Log().Infof(ctx, "文档保存成功, ID: %d", id)
	return id, nil
}

// UpdateDocumentsStatus 更新文档状态
func UpdateDocumentsStatus(ctx context.Context, documentsId int64, status int) error {
	data := g.Map{
		"status": status,
	}

	_, err := dao.KnowledgeDocuments.Ctx(ctx).Where("id", documentsId).Data(data).Update()
	if err != nil {
		g.Log().Errorf(ctx, "更新文档状态失败: ID=%d, 错误: %v", documentsId, err)
	}

	return err
}

// GetDocumentById 根据ID获取文档信息
func GetDocumentById(ctx context.Context, id int64) (document entity.KnowledgeDocuments, err error) {
	g.Log().Debugf(ctx, "获取文档信息: ID=%d", id)

	err = dao.KnowledgeDocuments.Ctx(ctx).Where("id", id).Scan(&document)
	if err != nil {
		g.Log().Errorf(ctx, "获取文档信息失败: ID=%d, 错误: %v", id, err)
		return document, fmt.Errorf("获取文档信息失败: %w", err)
	}

	return document, nil
}

// GetDocumentsList 获取文档列表
func GetDocumentsList(ctx context.Context, where entity.KnowledgeDocuments, page int, pageSize int) (documents []entity.KnowledgeDocuments, total int, err error) {
	// 参数验证和默认值设置
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	model := dao.KnowledgeDocuments.Ctx(ctx)
	if where.KnowledgeBaseName != "" {
		model = model.Where("knowledge_base_name", where.KnowledgeBaseName)
	}

	total, err = model.Count()
	if err != nil {
		g.Log().Errorf(ctx, "获取文档总数失败: %v", err)
		return nil, 0, fmt.Errorf("获取文档总数失败: %w", err)
	}

	if total == 0 {
		return nil, 0, nil
	}

	err = model.Page(page, pageSize).
		Order("created_at desc").
		Scan(&documents)
	if err != nil {
		g.Log().Errorf(ctx, "获取文档列表失败: %v", err)
		return nil, 0, fmt.Errorf("获取文档列表失败: %w", err)
	}

	return documents, total, nil
}

// DeleteDocument 删除文档及其相关数据
func DeleteDocument(ctx context.Context, id int64) error {
	g.Log().Debugf(ctx, "删除文档: ID=%d", id)

	return dao.KnowledgeDocuments.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先删除文档块
		_, err := dao.KnowledgeChunks.Ctx(ctx).TX(tx).Where("knowledge_doc_id", id).Delete()
		if err != nil {
			g.Log().Errorf(ctx, "删除文档块失败: ID=%d, 错误: %v", id, err)
			return fmt.Errorf("删除文档块失败: %w", err)
		}

		// 再删除文档
		result, err := dao.KnowledgeDocuments.Ctx(ctx).TX(tx).Where("id", id).Delete()
		if err != nil {
			g.Log().Errorf(ctx, "删除文档失败: ID=%d, 错误: %v", id, err)
			return fmt.Errorf("删除文档失败: %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("获取影响行数失败: %w", err)
		}
		if affected == 0 {
			return fmt.Errorf("文档不存在")
		}

		g.Log().Infof(ctx, "文档删除成功: ID=%d", id)
		return nil
	})
}
