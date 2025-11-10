package knowledge

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/internal/dao"
	"github.com/wangle201210/go-rag/server/internal/model/entity"
)

// SaveChunksData 批量保存知识块数据
func SaveChunksData(ctx context.Context, documentsId int64, chunks []entity.KnowledgeChunks) error {
	if len(chunks) == 0 {
		return nil
	}
	status := int(v1.StatusIndexing)
	// 使用 OnConflict 指定冲突列（chunk_id 是唯一索引）
	// 当 chunk_id 冲突时，更新其他字段
	_, err := dao.KnowledgeChunks.Ctx(ctx).Data(chunks).
		OnConflict("chunk_id").
		Save()
	if err != nil {
		g.Log().Errorf(ctx, "SaveChunksData err=%+v", err)
		status = int(v1.StatusFailed)
	}
	UpdateDocumentsStatus(ctx, documentsId, status)
	return err
}

// GetChunksList 查询知识块列表
func GetChunksList(ctx context.Context, where entity.KnowledgeChunks, page, size int) (list []entity.KnowledgeChunks, total int, err error) {
	model := dao.KnowledgeChunks.Ctx(ctx)

	// 构建查询条件
	if where.KnowledgeDocId != 0 {
		model = model.Where("knowledge_doc_id", where.KnowledgeDocId)
	}
	if where.ChunkId != "" {
		model = model.Where("chunk_id", where.ChunkId)
	}

	// 获取总数
	total, err = model.Count()
	if err != nil {
		return
	}

	// 分页查询
	if page > 0 && size > 0 {
		model = model.Page(page, size)
	}

	// 按创建时间倒序
	model = model.OrderDesc("created_at")

	err = model.Scan(&list)
	return
}

// GetChunkById 根据ID查询单个知识块
func GetChunkById(ctx context.Context, id int64) (chunk entity.KnowledgeChunks, err error) {
	err = dao.KnowledgeChunks.Ctx(ctx).Where("id", id).Scan(&chunk)
	return
}

// DeleteChunkByIds 根据ID软删除知识块
func DeleteChunkById(ctx context.Context, id int64) error {
	_, err := dao.KnowledgeChunks.Ctx(ctx).Where("id", id).Delete()
	return err
}

// UpdateChunkById 根据ID更新知识块
func UpdateChunkByIds(ctx context.Context, ids []int64, data entity.KnowledgeChunks) error {
	model := dao.KnowledgeChunks.Ctx(ctx).WhereIn("id", ids)
	if data.Content != "" {
		model = model.Data("content", data.Content)
	}
	if data.Status != 0 {
		model = model.Data("status", data.Status)
	}
	_, err := model.Update()
	return err
}

// GetAllChunksByDocId gets all chunks by document id
func GetAllChunksByDocId(ctx context.Context, docId int64, fields ...string) (list []entity.KnowledgeChunks, err error) {
	model := dao.KnowledgeChunks.Ctx(ctx).Where("knowledge_doc_id", docId)
	if len(fields) > 0 {
		for _, field := range fields {
			model = model.Fields(field)
		}
	}
	err = model.Scan(&list)
	return
}
