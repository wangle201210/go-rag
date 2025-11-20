package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

type SQLite struct {
	cfg *Config
	Base
}

func NewSqlite(cfg *Config) *SQLite {
	s := &SQLite{cfg: cfg}
	return s
}

func (s *SQLite) Connect() error {
	dsn := s.DSN()
	dialect := sqlite.Open(dsn)
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(s.cfg.LogLevel)),
	}
	db, err := gorm.Open(dialect, config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(s.cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(s.cfg.MaxOpenConn)

	sqlDB.SetConnMaxLifetime(time.Hour)

	s.db = db
	return nil
}

func (s *SQLite) DSN() string {
	var dsn string
	var params []string
	if s.cfg.JournalNode != "" {
		params = append(params, fmt.Sprintf("_journal_mode=%s", s.cfg.JournalNode))
	}
	if s.cfg.Synchronous != "" {
		params = append(params, fmt.Sprintf("_synchronous=%s", s.cfg.Synchronous))
	}
	if s.cfg.CacheSize > 0 {
		params = append(params, fmt.Sprintf("_cache_size=%d", s.cfg.CacheSize))
	}
	if s.cfg.BusyTimeout > 0 {
		params = append(params, fmt.Sprintf("_busy_timeout=%d", s.cfg.BusyTimeout))
	}
	dsn = s.cfg.FilePath
	if len(params) > 0 {
		dsn += "?" + strings.Join(params, "&")
	}
	return dsn
}

