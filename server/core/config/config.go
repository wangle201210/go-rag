package config

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Client    *elasticsearch.Client
	IndexName string // es index name
	// embedding 时使用
	APIKey  string
	BaseURL string
	Model   string
}
