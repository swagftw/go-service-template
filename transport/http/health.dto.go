package http

type PingReq struct {
	Name string `json:"name" validate:"required,min=3"`
}

type PingRes struct {
	GreetMessage string `json:"greetMessage"`
}
