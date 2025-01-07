package logrepository

type LogRepository interface {
	AddLogEntry(string, string) error
}
