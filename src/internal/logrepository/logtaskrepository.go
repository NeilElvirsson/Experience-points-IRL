package logrepository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type logTaskRepository struct {
	databasePath string
}

func New(dbPath string) logTaskRepository {

	return logTaskRepository{
		databasePath: dbPath,
	}
}

func (ltrepo logTaskRepository) AddLogEntry(userId string, taskId string) error {

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", ltrepo.databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO log(id, user_id, task_id, timestamp) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id.String(), userId, taskId, time.Now().Unix())
	if err != nil {
		return err
	}

	value, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if value == 0 {
		return errors.New("failed to add log")
	}
	return nil

}
