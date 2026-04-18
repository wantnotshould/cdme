// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package database

import (
	"log"
	"time"

	"code.cn/blog/conf"
	"code.cn/blog/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DB struct {
	db *gorm.DB
}

var instance *DB

func newDB() (*DB, error) {
	cfg := conf.Get().Database
	dsn := cfg.DSN()

	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	maxIdle := cfg.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 10
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 100
	}

	maxLife := cfg.MaxLifetime
	if maxLife == 0 {
		maxLife = time.Hour
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLife)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &DB{db: gdb}, nil
}

func Instance() *DB {
	return instance
}

func Init() {
	var err error
	instance, err = newDB()
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
}

func Get() *gorm.DB {
	if instance == nil {
		return nil
	}
	return instance.db
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserToken{},
		&model.Post{},
	)
}

func (d *DB) Close() error {
	if d == nil || d.db == nil {
		return nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
