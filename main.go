package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/numanijaz/tinyurl/config"
	"github.com/numanijaz/tinyurl/database"
	"github.com/numanijaz/tinyurl/handlers"
	"github.com/numanijaz/tinyurl/routers"
)

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func setupRoutersAndServe() {
	var mode string
	if config.GetConfig().GO_ENV == "production" {
		mode = gin.ReleaseMode
	} else {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	r := gin.Default()

	r.GET("/hearbeat", heartbeat)
	r.Static("/static", "./static")
	r.Static("/assets", "./build/assets")

	api := r.Group("/api")
	routers.SetupURLRoutes(api)
	routers.SetupAuthRouters(api)
	routers.SetupFrontendAppRouters(r)

	r.LoadHTMLFiles("./build/index.html")

	// serve tinyurl request
	r.GET("/:tinyurl", handlers.GetTinyUrl)

	cfg := config.GetConfig()
	r.Run(fmt.Sprintf("%s:%s", cfg.HOST_NAME, cfg.PORT))
}

func main() {
	config.GetConfig() // load the config
	config.SetupOauthProviders()

	database.InitAndMigrateDB()
	setupRoutersAndServe()
}
