package logrepository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/NeilElvirsson/Experience-points-IRL/internal/models"
	"github.com/google/uuid"
)

//
//

var ErrLogNotFound = errors.New("log not found")

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

func (ltrepo logTaskRepository) GetLogs(userId string) ([]models.Log, error) {

	db, err := sql.Open("sqlite3", ltrepo.databasePath)
	if err != nil {
		return []models.Log{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`
	SELECT 
		l.task_id, l.timestamp, t.task_name, t.xp_value 
	FROM 
		log l
	INNER JOIN 
		task t ON l.task_id = t.id
	WHERE l.user_id = ?`)
	if err != nil {
		return []models.Log{}, err

	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return []models.Log{}, err
	}
	defer rows.Close()

	var logs []models.Log

	for rows.Next() {
		var tempTaskId string
		var tempTimeStamp int
		var tempTaskName string
		var tempXpValue int

		err = rows.Scan(&tempTaskId, &tempTimeStamp, &tempTaskName, &tempXpValue)
		if err != nil {
			return []models.Log{}, err
		}

		logs = append(logs, models.Log{
			TaskId:    tempTaskId,
			Timestamp: tempTimeStamp,
			TaskName:  tempTaskName,
			XpValue:   tempXpValue,
		})

	}

	return logs, nil
}

func (ltrepo logTaskRepository) GetXpLevel(userId string) (models.XpSummary, error) {

	db, err := sql.Open("sqlite3", ltrepo.databasePath)
	if err != nil {
		return models.XpSummary{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`
	SELECT t.xp_value
		FROM log l
	INNER JOIN 
		task t ON l.task_id = t.id
	WHERE l.user_id = ?`)

	if err != nil {
		return models.XpSummary{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return models.XpSummary{}, err
	}
	defer rows.Close()

	var XpSummary models.XpSummary
	totalXp := 0
	levelDecimal := 1.0
	baseLevelXp := 100.0
	growthRate := 1.25
	progress := 0.0

	for rows.Next() {
		var tempXpValue int

		err = rows.Scan(&tempXpValue)
		if err != nil {
			return models.XpSummary{}, err
		}
		totalXp = totalXp + tempXpValue

	}
	XpSummary.TotalXp = totalXp
	remainingXp := float64(totalXp)

	for {
		xpForNextLevel := baseLevelXp * math.Pow(growthRate, levelDecimal-1)

		if remainingXp < xpForNextLevel {
			progress = (remainingXp / xpForNextLevel) * 100
			XpSummary.Progress = int(progress)
			fmt.Println("remainingXP: ", remainingXp, "xpForNextLevel: ", xpForNextLevel)
			break
		}

		remainingXp = remainingXp - xpForNextLevel
		levelDecimal++

	}
	level := int(levelDecimal)

	XpSummary.Level = level

	return XpSummary, nil

}
