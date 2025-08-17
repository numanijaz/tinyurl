package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/numanijaz/tinyurl/handlers"
	"github.com/numanijaz/tinyurl/routers"
)

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func setupRouters() {
	r := gin.Default()

	r.GET("/hearbeat", heartbeat)
	r.Static("/static", "./static")
	r.Static("/assets", "./build/assets")

	api := r.Group("/api")
	routers.SetupURLRoutes(api)
	routers.SetupAuthRouters(api)

	r.LoadHTMLFiles("./build/index.html")
	r.GET("/:tinyurl", handlers.GetTinyUrl)
	r.GET("/", handlers.ServeFrontendApp)
	r.GET("/notfound", handlers.ServeFrontendApp)

	r.Run(":8000")
}

func main() {

	// InitDB()
	// DB.AutoMigrate(&models.UrlModel{})

	setupRouters()
}
