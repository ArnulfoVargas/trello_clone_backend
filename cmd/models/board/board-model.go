package models_board

import (
	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models"
	models_auth "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/auth"
)

type Board struct {
	models.Base
	Name   string
	UserID string
	User   models_auth.User
}

type BoardIn struct {
	Name string
}

type BoardOut struct {
	ID string
	Name string
}

