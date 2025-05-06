package models

type User struct {
	UserName string
	UserId   string
	Password string
}

type Log struct {
	TaskId    string
	Timestamp int
	TaskName  string
	XpValue   int
}

type XpSummary struct {
	TotalXp  int
	Level    int
	Progress int
}
