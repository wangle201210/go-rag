// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KnowledgeBase is the golang structure of table knowledge_base for DAO operations like Where/Data.
type KnowledgeBase struct {
	g.Meta      `orm:"table:knowledge_base, do:true"`
	Id          interface{} // 主键ID
	Name        interface{} // 知识库名称
	Description interface{} // 知识库描述
	Category    interface{} // 知识库分类
	Status      interface{} // 状态：0-禁用，1-启用
	CreateTime  *gtime.Time // 创建时间
	UpdateTime  *gtime.Time // 更新时间
}
