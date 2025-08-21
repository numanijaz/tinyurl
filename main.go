package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/numanijaz/tinyurl/config"
	"github.com/numanijaz/tinyurl/handlers"
	"github.com/numanijaz/tinyurl/routers"
)

func heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func setupOauthProviders() {
	cfg := config.GetConfig()
	// gothic
	goth.UseProviders(
		github.New(
			cfg.GITHUB_CLIENT_ID,
			cfg.GITHUB_CLIENT_SECRET,
			cfg.BASE_URL+"/api/auth/callback/github",
			"user:email",
		),
		google.New(
			cfg.GOOGLE_CLIENT_ID,
			cfg.GOOGLE_CLIENT_SECRET,
			cfg.BASE_URL+"/api/auth/callback/google",
			"email",
			"profile",
		),
	)
	// gothic.Store = config.CookieStore
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
	r.GET("/error", handlers.ServeFrontendApp)
	r.GET("/notfound", handlers.ServeFrontendApp)

	cfg := config.GetConfig()
	r.Run(fmt.Sprintf("%s:%s", cfg.HOST_NAME, cfg.PORT))
}

func main() {
	config.GetConfig() // load the config
	setupOauthProviders()

	// InitDB()
	// DB.AutoMigrate(&models.UrlModel{})

	setupRouters()
}
