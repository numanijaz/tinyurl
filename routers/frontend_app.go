package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/numanijaz/tinyurl/handlers"
)

var frontendRoutes = []string{
	"/",
	"/error",
	"/notfound",
	"/login",
	"/login/*any",
	"/signup",
	"/signup/*any",
}

func SetupFrontendAppRouters(engine *gin.Engine) {
	for _, v := range frontendRoutes {
		engine.GET(v, handlers.ServeFrontendApp)
	}
}
