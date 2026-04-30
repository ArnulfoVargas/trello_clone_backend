package models_auth

import (
	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models"
)

type User struct {
	models.Base
	Name string `gorm:"type:varchar(50);not null"`
	Email string `gorm:"type:varchar(150);uniqueIndex;not null"`
	Password string `gorm:"type:char(60);not null"`
	Role uint8 `gorm:"default:0"`
}
