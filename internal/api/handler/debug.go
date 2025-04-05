package handler

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/vcsfrl/xm/internal/api/middleware"
	"net/http"
)

func NewDebug(logger zerolog.Logger) http.Handler {
	//TODO: make this work with gin
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.LoggerMiddleware(logger))
	r.Mount("/debug", chiMiddleware.Profiler())

	return r
}
