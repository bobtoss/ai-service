package router

import (
	_ "ai-service/docs"
	"ai-service/internal/repository"
	"ai-service/internal/repository/postgres"
	"ai-service/internal/service/auth"
	"ai-service/internal/service/chat"
	"ai-service/internal/service/doc"
	"ai-service/internal/util/config"
	authMiddleware "ai-service/internal/util/middleware"
	"ai-service/internal/util/validator"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	config     *config.Config
	repository *repository.Repository
}

func NewRouter(cfg *config.Config, repo *repository.Repository) *Router {
	return &Router{
		config:     cfg,
		repository: repo,
	}
}

func (r Router) Build(ctx context.Context) *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()
	e.Pre(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	db, err := postgres.NewDB(ctx, r.config.DB.DSN())
	if err != nil {
		panic(err)
	}
	userRepo := postgres.NewUserRepository(db)
	jwtSecret := []byte(r.config.JWTSecret)

	docRepo := postgres.NewDocumentRepository(db)
	docService, err := doc.NewDocService(r.config, r.repository, docRepo)
	if err != nil {
		panic(err)
	}

	svc := auth.NewService(userRepo, jwtSecret)
	authHandler := auth.NewAuthHandler(svc)
	authMw := authMiddleware.AuthMiddleware(svc.ParseAccessToken)

	{
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.POST("/login", authHandler.Login)
		e.POST("/refresh", authHandler.Refresh)
		e.POST("/logout", authHandler.Logout)
		e.POST("/register", authHandler.Register)
	}
	api := e.Group("", authMw)
	{
		services := api.Group("/upload")
		services.POST("", docService.SaveDoc)
		services.GET("", docService.ListDoc)
		services.PUT("/:id", docService.UpdatePriority)
		services.DELETE("/:id", docService.DeleteDoc)
	}
	chatService, err := chat.NewChatService(r.config, r.repository)
	if err != nil {
		panic(err)
	}
	{
		services := api.Group("/chat")
		services.POST("", chatService.Chat)
	}
	return e
}
