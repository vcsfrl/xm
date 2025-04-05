package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/vcsfrl/xm/internal/api/handler"
	"github.com/vcsfrl/xm/internal/api/middleware"
	"github.com/vcsfrl/xm/internal/api/validator"
	"github.com/vcsfrl/xm/internal/config"
	"github.com/vcsfrl/xm/internal/service"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const serverShutdownDelay = 5

type RestApi struct {
	ctx    context.Context
	logger zerolog.Logger
	config *config.Config
	db     *gorm.DB
	srv    *http.Server
}

func NewRestApi(ctx context.Context, logger zerolog.Logger, config *config.Config, db *gorm.DB) *RestApi {
	return &RestApi{ctx: ctx, logger: logger, config: config, db: db}
}

func (c *RestApi) Run() {
	c.logger.Info().Msg("Running api")

	router, err := c.BuildRouter()
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to build router")
		return
	}

	c.srv = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", c.config.AppPort),
		Handler: router,
	}

	if err := c.srv.ListenAndServe(); err != nil {
		c.logger.Error().Err(err).Msg("Failed to run the server")
		return
	}
}

func (c *RestApi) BuildRouter() (*gin.Engine, error) {
	companyService := service.NewCompanyService(c.db, validator.CompanyValidator(c.logger))
	companyHandler := handler.NewCompanyHandler(companyService)

	authManager, err := middleware.NewAuthenticationManager(c.config, c.logger)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to create auth manager")
		return nil, err
	}

	ginRouter := gin.Default()
	ginRouter.Use(middleware.RateLimiter(c.config))
	ginRouter.Use(authManager.JwtHandler())
	apiRouter := ginRouter.Group("/api/v1")
	apiRouter.POST("/login", authManager.AuthMiddleware.LoginHandler)
	apiRouter.GET("/health", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	apiRouter.GET("/company/:id", companyHandler.Get)

	// register middleware
	authorized := apiRouter.Group("/", authManager.AuthMiddleware.MiddlewareFunc())
	{
		authorized.POST("/company", companyHandler.Create)
		authorized.PATCH("/company/:id", companyHandler.Update)
		authorized.DELETE("/company/:id", companyHandler.Delete)

		authorized.POST("/refresh_token", authManager.AuthMiddleware.RefreshHandler)
	}

	return ginRouter, nil
}

func (c *RestApi) Close() error {
	c.logger.Info().Msg("Shutting down api")
	if c.srv == nil {
		c.logger.Info().Msg("Server is not running")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownDelay*time.Second)
	defer cancel()
	if err := c.srv.Shutdown(ctx); err != nil {
		c.logger.Error().Err(err).Msg("Failed to shutdown server")
		return err
	}

	c.logger.Info().Msg("Server shutdown successfully")
	return nil
}
