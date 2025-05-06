package logrepository

import "github.com/NeilElvirsson/Experience-points-IRL/internal/models"

type LogRepository interface {
	AddLogEntry(string, string) error
	GetLogs(string) ([]models.Log, error)
	GetXpLevel(string) (models.XpSummary, error)
}
