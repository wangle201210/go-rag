package knowledge_base

import (
	"context"

	"github.com/wangle201210/go-rag/server/internal/model"

	"github.com/wangle201210/go-rag/server/api/knowledge_base/v1"
)

func (c *ControllerV1) KnowledgeBase(ctx context.Context, req *v1.KnowledgeBaseReq) (res *v1.KnowledgeBaseRes, err error) {
	list, err := c.knowledgeBaseService.GetList(ctx, &model.KnowledgeBaseListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Category: req.Category,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return
	}
	res = &v1.KnowledgeBaseRes{
		List:  model2Data(list.List),
		Total: list.Total,
	}
	return
}

func model2Data(data []model.KnowledgeBaseItem) (res []v1.KnowledgeBaseItem) {
	for _, item := range data {
		res = append(res, v1.KnowledgeBaseItem{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Category:    item.Category,
			Status:      item.Status,
			CreateTime:  item.CreateTime,
			UpdateTime:  item.UpdateTime,
		})
	}
	return
}
