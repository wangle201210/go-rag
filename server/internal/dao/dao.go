package dao

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	gormModel "github.com/wangle201210/go-rag/server/internal/model/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

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
	if dbType == "" {
		dbType = "mysql" // 默认使用 mysql
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var dialector gorm.Dialector
	var err error

	switch dbType {
	case "sqlite":
		filePath := g.Cfg().MustGet(ctx, "database.default.host").String()
		if filePath == "" {
			return fmt.Errorf("sqlite file path is required when using sqlite")
		}
		// 确保 SQLite 文件的父目录存在
		if err := ensureSQLiteFileDir(filePath); err != nil {
			return fmt.Errorf("failed to create sqlite file directory: %v", err)
		}
		dsn := buildSQLiteDSN(ctx, filePath)
		dialector = sqlite.Open(dsn)
	case "mysql":
		dsn := GetDsn()
		dialector = mysql.Open(dsn)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	db, err = gorm.Open(dialector, config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池
	if dbType == "sqlite" {
		// SQLite 建议使用单连接
		maxOpenConns := g.Cfg().MustGet(ctx, "database.default.max_open_conns", 1).Int()
		maxIdleConns := g.Cfg().MustGet(ctx, "database.default.max_idle_conns", 1).Int()
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表结构
	if err = gormModel.AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate database tables: %v", err)
	}

	return nil
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

// buildSQLiteDSN 构建 SQLite DSN
func buildSQLiteDSN(ctx context.Context, filePath string) string {
	dsn := filePath + "?"

	// 添加可选参数
	if busyTimeout := g.Cfg().MustGet(ctx, "database.default.busy_timeout").Int(); busyTimeout > 0 {
		dsn += fmt.Sprintf("_busy_timeout=%d&", busyTimeout)
	}
	if journalMode := g.Cfg().MustGet(ctx, "database.default.journal_mode").String(); journalMode != "" {
		dsn += fmt.Sprintf("_journal_mode=%s&", journalMode)
	}
	if synchronous := g.Cfg().MustGet(ctx, "database.default.synchronous").String(); synchronous != "" {
		dsn += fmt.Sprintf("_synchronous=%s&", synchronous)
	}
	if cacheSize := g.Cfg().MustGet(ctx, "database.default.cache_size").Int(); cacheSize != 0 {
		dsn += fmt.Sprintf("_cache_size=%d&", cacheSize)
	}

	// 移除末尾的 & 或 ?
	if dsn[len(dsn)-1] == '&' || dsn[len(dsn)-1] == '?' {
		dsn = dsn[:len(dsn)-1]
	}

	return dsn
}

func GetDsn() string {
	ctx := context.Background()
	dbType := GetDBType()

	switch dbType {
	case "sqlite":
		// 对于 SQLite，返回简单的文件路径格式
		// chat-history 包会自动处理 SQLite 的 DSN
		filePath := g.Cfg().MustGet(ctx, "database.default.file").String()
		// 如果有 journal_mode 配置，添加到 DSN
		if journalMode := g.Cfg().MustGet(ctx, "database.default.journal_mode").String(); journalMode != "" {
			return fmt.Sprintf("%s?_journal_mode=%s", filePath, journalMode)
		}
		return filePath
	case "mysql":
		cfg := g.DB().GetConfig()
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.Charset)
	default:
		g.Log().Fatalf(ctx, "unsupported database type: %s", dbType)
		return ""
	}
}

// GetDBType 获取数据库类型
func GetDBType() string {
	ctx := context.Background()
	dbType := g.Cfg().MustGet(ctx, "database.default.type", "mysql").String()
	if dbType == "" {
		return "mysql"
	}
	return dbType
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		g.Log().Fatal(context.Background(), "database connection not initialized")
	}
	return db
}
