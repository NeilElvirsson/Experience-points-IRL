package taskrepository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Package taskrepository handles all interactions with the database related to tasks.

type addTaskRepository struct {
	databasePath string
}

func New(dbPath string) addTaskRepository {
	return addTaskRepository{
		databasePath: dbPath,
	}
}

func (atrepo addTaskRepository) AddTask(taskName string, xpValue int) error {

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", atrepo.databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO task (id, task_name, xp_value) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id.String(), taskName, xpValue)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {

		fmt.Println("failed to add task")
	}
	fmt.Println("Task Id: ", id)
	return nil

}
