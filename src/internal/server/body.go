package server

type addUserRequestBody struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
