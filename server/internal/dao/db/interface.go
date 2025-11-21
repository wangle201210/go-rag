package db

// Database 数据库接口
type Database interface {
	// Connect 连接管理
	Connect() error
	// DSN 构建连接DSN
	DSN() string
	// AutoMigrate 迁移表结构
	AutoMigrate() error
	// Ping 健康检查
	Ping() error
}

type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	Charset     string
	MaxOpenConn int
	MaxIdleConn int
	LogLevel    int // 1: silent 2: Error 3:Warn 4: Info

	// sqlite相关配置
	FilePath    string
	BusyTimeout int
	JournalMode string
	Synchronous string
	CacheSize   int
}
