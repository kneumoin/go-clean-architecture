package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kneumoin/go-clean-architecture/link"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc link.UseCase) {
	h := NewHandler(uc)

	links := router.Group("/links")
	{
		links.POST("", h.Create)
		links.GET("/:id", h.Get)
	}
}
