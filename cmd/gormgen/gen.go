package main

import (
	"go-service-template/infrastructure/db/postgres"
	"go-service-template/infrastructure/db/postgres/models"
	"go-service-template/internal/applog"
	"go-service-template/internal/config"
	"gorm.io/gen"
)

// GenerateTypeSafeQuery generate type safe query
func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		return
	}

	err = applog.InitLogger(cfg.Debug)
	if err != nil {
		return
	}

	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		return
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./infrastructure/db/postgres/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db) // reuse your gorm db

	g.ApplyBasic(models.Ping{})

	// Generate the code
	g.Execute()
}
