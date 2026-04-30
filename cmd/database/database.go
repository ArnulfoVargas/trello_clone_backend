package database

import (
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

	db.Db = conn

	db.bindMigrations()

	Connection = db
}
