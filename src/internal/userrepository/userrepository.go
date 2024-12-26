package userrepository

import (
	"github.com/NeilElvirsson/Experience-points-IRL/internal/models"
)

type UserRepository interface {
	LoginUser(string, string) (models.User, error)
	AddUser(models.User) error
}
