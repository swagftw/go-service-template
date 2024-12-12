package postgres

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-service-template/internal/apperr"
	"go-service-template/internal/applog"
)

func Connect(url string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
		PrepareStmt:            true,
	})
	if err != nil {
		applog.Logger.Error(context.Background(), err, "failed to initialize database", map[string]interface{}{
			"url": url,
		})

		return nil, apperr.New(http.StatusInternalServerError, err, "failed to initialize database", apperr.ErrInternalError)
	}

	sqlDB, err := db.DB()
	if err != nil {
		applog.Logger.Error(context.Background(), err, "failed to initialize database", map[string]interface{}{
			"url": url,
		})

		return nil, apperr.New(http.StatusInternalServerError, err, "failed to initialize database", apperr.ErrInternalError)
	}

	err = sqlDB.Ping()
	if err != nil {
		applog.Logger.Error(context.Background(), err, "failed to initialize database", map[string]interface{}{
			"url": url,
		})

		return nil, apperr.New(http.StatusInternalServerError, err, "failed to initialize database", apperr.ErrInternalError)
	}

	return db, nil
}
