package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/numanijaz/tinyurl/config"
)

func CurrentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		var err error
		if tokenStr, err = c.Cookie("authToken"); err != nil {
			c.Next()
			return
		}

		cfg := config.GetConfig()
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return []byte(cfg.SECRET_KEY), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			c.Next()
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			c.Next()
			return
		}

		c.Set("sub", sub)
	}
}
