package db

import (
	"github.com/vcsfrl/xm/internal/config"
	"github.com/vcsfrl/xm/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqlite(config *config.Config) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	// TODO integrate zerolog
	//gormLogger := gormzerolog.NewGormLogger().WithInfo(func() gormzerolog.Event {
	//	return &gormzerolog.GormLoggerEvent{Event: logger.Info()}
	//})

	db, err = gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Company{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitTestSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.Company{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
