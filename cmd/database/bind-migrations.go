package database

import models_auth "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/auth"

func (db *Database) bindMigrations() {
	db.Db.AutoMigrate(&models_auth.User{})
}
