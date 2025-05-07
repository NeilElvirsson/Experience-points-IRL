package userrepository

import (
	"database/sql"
	"errors"

	"github.com/NeilElvirsson/Experience-points-IRL/internal/models"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Package userrepository provides functionality for managing user data in the database.
// It supports creating new users and authenticating existing ones using SQLite.

var ErrUserNotFound = errors.New("user not found")

type sqliteUserRepository struct {
	databasePath string
}

func New(dbPath string) sqliteUserRepository {

	return sqliteUserRepository{
		databasePath: dbPath,
	}
}

func (sqlite sqliteUserRepository) LoginUser(userName string, password string) (models.User, error) {

	db, err := sql.Open("sqlite3", sqlite.databasePath)
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, user_name FROM user WHERE user_name = ? AND password = ?")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userName, password)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var tempUserId string
		var tempUserName string

		err := rows.Scan(&tempUserId, &tempUserName)
		if err != nil {
			return models.User{}, err
		}

		return models.User{
			UserName: tempUserName,
			UserId:   tempUserId,
			Password: password,
		}, nil

	}
	return models.User{}, ErrUserNotFound

}

func (sqlite sqliteUserRepository) AddUser(user models.User) error {

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", sqlite.databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user(id, user_name, password) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id.String(), user.UserName, user.Password)
	if err != nil {
		return err
	}

	value, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if value == 0 {
		return errors.New("failed to add user")
	}

	return nil

}
