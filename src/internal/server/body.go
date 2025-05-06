package server

type addUserRequestBody struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type addLogTaskRequestBody struct {
	TaskId string `json:"taskId"`
}

type addTaskRequestBody struct {
	TaskName string `json:"taskName"`
	XpValue  int    `json:"xpValue"`
}

type getResponseBody struct {
	TaskId    string `json:"taskId"`
	Timestamp int    `json:"timestamp"`
	TaskName  string `json:"taskName"`
	XpValue   int    `json:"xpValue"`
}

type getXpSummaryBody struct {
	TotalXp  int `json:"totalXp"`
	Level    int `json:"level"`
	Progress int `json:"progress %"`
}
