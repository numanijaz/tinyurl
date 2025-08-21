package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/numanijaz/tinyurl/handlers"
)

func SetupURLRoutes(rg *gin.RouterGroup) {
	rg.POST("shortenurl", handlers.ShortenUrl)
}
