package common

import "github.com/cloudwego/eino/schema"

type StreamData struct {
	Id       string             `json:"id"`      // 同一个消息里面的id是相同的
	Created  int64              `json:"created"` // 消息初始生成时间
	Content  string             `json:"content"` // 消息具体内容
	Document []*schema.Document `json:"document"`
}
