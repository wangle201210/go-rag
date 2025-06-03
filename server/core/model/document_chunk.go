package model

import "time"

type DocumentChunk struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	DocumentId string    `json:"document_id" gorm:"index"`
	ChunkId    string    `json:"chunk_id" gorm:"uniqueIndex"`
	Content    string    `json:"content" gorm:"type:text"`
	Metadata   string    `json:"metadata" gorm:"type:text"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type DocumentChunkCreateReq struct {
	DocumentId string `json:"document_id" validate:"required"`
	ChunkId    string `json:"chunk_id" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Metadata   string `json:"metadata" validate:"required"`
}

type DocumentChunkCreateRes struct {
	Id int64 `json:"id"`
}

type DocumentChunkUpdateReq struct {
	Id       int64  `json:"id" validate:"required"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}

type DocumentChunkUpdateRes struct {
	Id int64 `json:"id"`
}

type DocumentChunkDeleteReq struct {
	Id int64 `json:"id" validate:"required"`
}

type DocumentChunkDeleteRes struct {
	Id int64 `json:"id"`
}

type DocumentChunkListReq struct {
	DocumentId string `json:"document_id" validate:"required"`
	Page       int    `json:"page" validate:"required,min=1"`
	PageSize   int    `json:"page_size" validate:"required,min=1,max=100"`
}

type DocumentChunkListRes struct {
	Total int64           `json:"total"`
	List  []DocumentChunk `json:"list"`
}
