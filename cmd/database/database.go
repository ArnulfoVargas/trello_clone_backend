package database

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var config = gorm.Config{
	DisableForeignKeyConstraintWhenMigrating: true,
}

var Connection *Database

type Database struct {
	Db *gorm.DB
}

func (db *Database) ConnectDB(dsn string) {
	if Connection != nil {
		return
	}

	conn, err := gorm.Open(mysql.Open(dsn), &config)

	if err != nil {
		panic(err)
	}

	sqlDB, err := conn.DB()
    if err != nil {
        panic(err)
    }

    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	go func() {
        for {
            time.Sleep(time.Minute * 5)
            sqlDB.Ping()
        }
    }()

	db.Db = conn

	db.bindMigrations()

	Connection = db
}
