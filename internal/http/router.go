package api

import (
	"net/http"

	"go-api-boilerplate/internal/config"
	"go-api-boilerplate/internal/modules/example"
	"go-api-boilerplate/internal/shared"

	"gorm.io/gorm"
)

func NewRouter(gormDB *gorm.DB, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteSuccess(w, http.StatusOK, "service is healthy", nil)
	})

	exampleRepo := example.NewRepository(gormDB)
	exampleService := example.NewService(exampleRepo)
	example.NewHandler(exampleService).Register(mux)

	handler := apply(mux,
		corsMiddleware(cfg.CORSAllowedOrigins),
		loggingMiddleware,
		apiKeyMiddleware(cfg.APIKey),
		authMiddleware(cfg.JWTSecret),
		permissionMiddleware(),
	)

	return handler
}
