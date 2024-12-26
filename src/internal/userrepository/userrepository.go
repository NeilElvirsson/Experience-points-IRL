package userrepository

import (
	"github.com/NeilElvirsson/Experience-points-IRL/internal/domain"
)

type UserRepository interface {
	LoginUser(string, string) (domain.User, error)
	AddUser(domain.User) error
}
