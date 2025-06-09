package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/internal/model/entity"
)

// Status marks kb status.
type Status int

const (
	StatusOK       Status = 1
	StatusDisabled Status = 2
)

type KBCreateReq struct {
	g.Meta      `path:"/v1/kb" method:"post" tags:"kb" summary:"Create kb"`
	Name        string `v:"required|length:3,20" dc:"kb name"`
	Description string `v:"required|length:3,200" dc:"kb description"`
	Category    string `v:"length:3,10" dc:"kb category"`
}

type KBCreateRes struct {
	Id int64 `json:"id" dc:"kb id"`
}

type KBUpdateReq struct {
	g.Meta      `path:"/v1/kb/{id}" method:"put" tags:"kb" summary:"Update kb"`
	Id          int64   `v:"required" dc:"kb id"`
	Name        *string `v:"length:3,10" dc:"kb name"`
	Description *string `v:"length:3,200" dc:"kb description"`
	Category    *string `v:"length:3,10" dc:"kb category"`
	Status      *Status `v:"in:1,2" dc:"kb status"`
}
type KBUpdateRes struct{}

type KBDeleteReq struct {
	g.Meta `path:"/v1/kb/{id}" method:"delete" tags:"kb" summary:"Delete kb"`
	Id     int64 `v:"required" dc:"kb id"`
}
type KBDeleteRes struct{}

type KBGetOneReq struct {
	g.Meta `path:"/v1/kb/{id}" method:"get" tags:"kb" summary:"Get one kb"`
	Id     int64 `v:"required" dc:"kb id"`
}
type KBGetOneRes struct {
	*entity.KnowledgeBase `dc:"kb"`
}

type KBGetListReq struct {
	g.Meta   `path:"/v1/kb" method:"get" tags:"kb" summary:"Get kbs"`
	Name     *string `v:"length:3,10" dc:"kb name"`
	Status   *Status `v:"in:1,2" dc:"kb age"`
	Category *string `v:"length:3,10" dc:"kb category"`
}

type KBGetListRes struct {
	List []*entity.KnowledgeBase `json:"list" dc:"kb list"`
}
