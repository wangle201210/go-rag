package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
	Base
	cfg *Config
}

func NewMysql(cfg *Config) *Mysql {
	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	return &Mysql{cfg: cfg}
}

func (m *Mysql) Connect() error {
	dsn := m.DSN()
	dialect := mysql.Open(dsn)
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(m.cfg.LogLevel)),
	}
	db, err := gorm.Open(dialect, config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(m.cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(m.cfg.MaxOpenConn)

	sqlDB.SetConnMaxLifetime(time.Hour)

	m.DB = db
	return nil
}

func (m *Mysql) DSN() string {
	var dsn string
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		m.cfg.User, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.Database, m.cfg.Charset)
	return dsn
}
