package router

import (
	"ai-service/internal/repository"
	"ai-service/internal/service/doc"
	"ai-service/internal/util/config"
	"github.com/labstack/echo/v4"
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

func (r Router) Build() *echo.Echo {
	e := echo.New()
	api := e.Group("/api/v1")
	docService, err := doc.NewDocService(r.config, r.repository)
	if err != nil {
		panic(err)
	}
	{
		services := api.Group("/doc")
		services.POST("", docService.SaveDoc)
		services.PUT("/:id", docService.UpdatePriority)
		services.DELETE(":id", docService.DeleteDoc)
	}
	return e
}
