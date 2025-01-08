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
	XpValue  int    `json:"xpValue"` //int här funkar?
}
