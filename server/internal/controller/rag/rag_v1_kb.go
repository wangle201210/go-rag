package rag

import (
	"context"

	"github.com/wangle201210/go-rag/server/internal/dao"
	"github.com/wangle201210/go-rag/server/internal/model/do"

	"github.com/wangle201210/go-rag/server/api/rag/v1"
)

func (c *ControllerV1) KBCreate(ctx context.Context, req *v1.KBCreateReq) (res *v1.KBCreateRes, err error) {
	insertId, err := dao.KnowledgeBase.Ctx(ctx).Data(do.KnowledgeBase{
		Name:        req.Name,
		Status:      v1.StatusOK,
		Description: req.Description,
		Category:    req.Category,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	res = &v1.KBCreateRes{
		Id: insertId,
	}
	return
}

func (c *ControllerV1) KBDelete(ctx context.Context, req *v1.KBDeleteReq) (res *v1.KBDeleteRes, err error) {
	_, err = dao.KnowledgeBase.Ctx(ctx).WherePri(req.Id).Delete()
	return
}

func (c *ControllerV1) KBGetList(ctx context.Context, req *v1.KBGetListReq) (res *v1.KBGetListRes, err error) {
	res = &v1.KBGetListRes{}
	err = dao.KnowledgeBase.Ctx(ctx).Where(do.KnowledgeBase{
		Status:   req.Status,
		Name:     req.Name,
		Category: req.Category,
	}).Scan(&res.List)
	return
}

func (c *ControllerV1) KBGetOne(ctx context.Context, req *v1.KBGetOneReq) (res *v1.KBGetOneRes, err error) {
	res = &v1.KBGetOneRes{}
	err = dao.KnowledgeBase.Ctx(ctx).WherePri(req.Id).Scan(&res.KnowledgeBase)
	return
}

func (c *ControllerV1) KBUpdate(ctx context.Context, req *v1.KBUpdateReq) (res *v1.KBUpdateRes, err error) {
	_, err = dao.KnowledgeBase.Ctx(ctx).Data(do.KnowledgeBase{
		Name:        req.Name,
		Status:      req.Status,
		Description: req.Description,
		Category:    req.Category,
	}).WherePri(req.Id).Update()
	return
}
