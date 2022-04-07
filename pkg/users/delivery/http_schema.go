package delivery

type LoginResp struct {
	Token string `json:"token"`
}

type LoginReq struct {
	Username string
	Password string
}
