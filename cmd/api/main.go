package main

import (
	"log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/interfaces/web"
	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/cache"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/security"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	security.InitJWT(cfg.JWTSecretKey)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open gorm.DB: %v", err)
	}
	initDB(db)

	cache := cache.NewRedisCache(cfg.CacheURL)

	logger, err := zap.NewDevelopment()
	// use this in production environment
	// logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create zap logger")
	}
	defer logger.Sync()

	app := app.New(cfg, db, cache, logger)
	defer app.Close()

	router := web.InitRouter(app)

	if err := router.RunTLS(cfg.Addr, cfg.CertPath, cfg.KeyPath); err != nil {
		app.Logger.Fatal("Failed to start http server",
			zap.String("addr", cfg.Addr),
			zap.String("cert", cfg.CertPath),
		)
	}
}

func initDB(db *gorm.DB) {
	db.Migrator().AutoMigrate(
		&model.User{},
		&model.Movie{},
		&model.Showtime{},
		&model.Reservation{},
		&model.Hall{},
		&model.Seat{},
		&model.ShowtimeSeat{},
	)
}
