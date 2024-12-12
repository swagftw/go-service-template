package healthrepo

import (
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"go-service-template/domain"
	"go-service-template/infrastructure/db/postgres/models"
	"go-service-template/infrastructure/db/postgres/query"
	"go-service-template/internal/apperr"
	"go-service-template/internal/applog"
)

type repo struct {
	db *gorm.DB
	q  *query.Query
}

// Repo interface for the repo is inside the repo itself as there is no need to expose it outside
// only service can call the repo and not any other services, they need to call it via the repo's service
type Repo interface {
	// inputs and outputs are not repo specific but the service specific
	// which makes service compatible to another kinds of repos (cache, file based, etc.)
	CreatePing(ctx context.Context, req *domain.PingReq) (*domain.PingRes, error)
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db, q: query.Use(db)}
}

func (r repo) CreatePing(ctx context.Context, req *domain.PingReq) (*domain.PingRes, error) {
	ping := &models.Ping{
		Name: req.Name,
	}

	err := r.q.WithContext(ctx).Ping.Create(ping)
	if err != nil {
		applog.Logger.Error(ctx, err, "failed to create ping", map[string]interface{}{
			"ping": ping,
		})

		return nil, apperr.New(http.StatusInternalServerError, err, "failed to create ping", apperr.ErrInternalError)
	}

	return &domain.PingRes{GreetMessage: fmt.Sprintf("Hello %s!", req.Name)}, nil
}
