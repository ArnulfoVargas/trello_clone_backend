package database

import (
	models_auth "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/auth"
	models_board "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/board"
)

func (db *Database) bindMigrations() {
	db.Db.AutoMigrate(&models_auth.User{}, &models_board.Board{})
}
