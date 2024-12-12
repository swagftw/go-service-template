package domain

import "context"

type IHealthService interface {
	Ping(ctx context.Context, req *PingReq) (*PingRes, error)
}

type PingReq struct {
	Name string
}

type PingRes struct {
	GreetMessage string
}
