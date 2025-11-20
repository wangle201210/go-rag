package dao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	database "github.com/wangle201210/go-rag/server/internal/dao/db"
	"os"
	"path/filepath"
)

var db database.Database

func init() {
	err := InitDB()
	if err != nil {
		g.Log().Fatal(context.Background(), "database connection not initialized, err %v", err)
	}
}

// InitDB 初始化数据库连接
func InitDB() error {
	ctx := context.Background()
	dbType := g.Cfg().MustGet(ctx, "database.default.type", "mysql").String()
	var (
		cfg *database.Config
		err error
	)
	switch dbType {
	case "mysql", "":
		cfg, _ = GetMysqlConfig()
		db = database.NewMysql(cfg)
	case "sqlite":
		cfg, err = GetSqliteConfig()
		if err != nil {
			return err
		}
		db = database.NewSqlite(cfg)
	}
	if err = db.Connect(); err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	// 迁移表结构
	return db.AutoMigrate()
}

// ensureSQLiteFileDir 确保 SQLite 文件的父目录存在
func ensureSQLiteFileDir(filePath string) error {
	// 获取文件的目录路径
	dir := filepath.Dir(filePath)

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建目录（包括所有必要的父目录）
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		g.Log().Infof(context.Background(), "Created directory for SQLite database: %s", dir)
	}

	return nil
}

func GetDsn() string {
	return db.DSN()
}

func GetSqliteConfig() (*database.Config, error) {
	ctx := context.Background()
	cfg := &database.Config{}
	cfg.FilePath = g.Cfg().MustGet(ctx, "database.default.host").String()
	if cfg.FilePath == "" {
		return nil, fmt.Errorf("sqlite file path is required when using sqlite")
	}
	if err := ensureSQLiteFileDir(cfg.FilePath); err != nil {
		return nil, fmt.Errorf("failed to create sqlite file directory: %v", err)
	}
	cfg.BusyTimeout = g.Cfg().MustGet(ctx, "database.default.busy_timeout").Int()
	cfg.JournalNode = g.Cfg().MustGet(ctx, "database.default.journal_mode").String()
	cfg.Synchronous = g.Cfg().MustGet(ctx, "database.default.synchronous").String()
	cfg.CacheSize = g.Cfg().MustGet(ctx, "database.default.cache_size").Int()
	cfg.MaxOpenConn = g.Cfg().MustGet(ctx, "database.default.max_open_conns", 1).Int()
	cfg.MaxIdleConn = g.Cfg().MustGet(ctx, "database.default.max_idle_conns", 1).Int()
	cfg.LogLevel = 4
	return cfg, nil
}

func GetMysqlConfig() (*database.Config, error) {
	cfg := g.DB().GetConfig()
	c := &database.Config{
		Host:        cfg.Host,
		Port:        cfg.Port,
		User:        cfg.User,
		Password:    cfg.Pass,
		Database:    cfg.Name,
		Charset:     cfg.Charset,
		MaxOpenConn: cfg.MaxOpenConnCount,
		MaxIdleConn: cfg.MaxIdleConnCount,
		LogLevel:    4,
	}
	if c.MaxIdleConn == 0 {
		c.MaxIdleConn = 10
	}
	if c.MaxOpenConn == 0 {
		c.MaxOpenConn = 100
	}
	return c, nil
}
