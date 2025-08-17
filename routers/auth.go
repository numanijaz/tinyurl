package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/numanijaz/tinyurl/handlers"
	"github.com/numanijaz/tinyurl/middleware"
)

func SetupAuthRouters(api *gin.RouterGroup) {
	api.GET("/auth/me", middleware.CurrentUserMiddleware(), handlers.GetCurrentUserInfo)
	api.POST("/auth/register", handlers.RegisterUser)
	api.POST("/auth/login", handlers.Login)
	api.POST("/auth/logout", handlers.Logout)
}
