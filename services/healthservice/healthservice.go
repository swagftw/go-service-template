package healthservice

import (
	"context"

	"go-service-template/domain"
	"go-service-template/services/healthservice/healthrepo"
)

type service struct {
	repo healthrepo.Repo
}

func NewService(repo healthrepo.Repo) domain.IHealthService {
	return &service{repo: repo}
}

func (s service) Ping(ctx context.Context, req *domain.PingReq) (*domain.PingRes, error) {
	return s.repo.CreatePing(ctx, req)
}
