package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"github.com/numanijaz/tinyurl/config"
	"github.com/numanijaz/tinyurl/database"
	"github.com/numanijaz/tinyurl/models"
	"golang.org/x/crypto/bcrypt"
)

func UserExists(userEmail string) bool {
	var user models.UserModel
	database.DB.First(&user, "email = ?", userEmail)
	return user.Email != ""
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

func setAuthTokenCookie(user models.UserModel, c *gin.Context) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": strconv.Itoa(int(user.ID)),
		},
	)
	cfg := config.GetConfig()
	signedToken, err := token.SignedString([]byte(cfg.SECRET_KEY))
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{"error": "An error occured while signing JWT"},
		)
		return
	}

	c.SetCookie(
		"authToken",
		signedToken,
		int((15 * time.Minute).Seconds()),
		"/",
		"", // current domain
		config.GetConfig().GO_ENV == "production", // use secure in production mode
		true,
	)
}

func RegisterUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Request does not contain necessary arguments"},
		)
		return
	}

	if UserExists(email) {
		c.JSON(
			http.StatusOK, // due to security? email might be harvested with error response
			gin.H{"error": "Username/Email not available"},
		)
		return
	}

	hashedPassword, _ := HashPassword(password)
	user := models.UserModel{
		Email:          email,
		Name:           email,
		HashedPassword: hashedPassword,
	}
	database.DB.Create(&user)

	setAuthTokenCookie(user, c)
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	genericError := gin.H{"error": "Email or password incorrect."} // secure

	if email == "" || password == "" {
		c.JSON(
			http.StatusBadRequest,
			genericError,
		)
		return
	}

	var user models.UserModel
	database.DB.First(&user, "email = ?", email)

	if user.Email == "" {
		// user does not exist
		c.JSON(
			http.StatusBadRequest,
			genericError,
		)
		return
	}

	if !VerifyPassword(user.HashedPassword, password) {
		c.JSON(
			http.StatusBadRequest,
			genericError,
		)
		return
	}

	setAuthTokenCookie(user, c)
	c.JSON(http.StatusOK, gin.H{})
}

func GetCurrentUserInfo(c *gin.Context) {
	subject, exists := c.Get("sub")
	if !exists {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	var user models.UserModel
	database.DB.First(&user, "id = ?", subject.(string))
	if user.Email == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"email": user.Email,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie(
		"authToken",
		"",
		-1,
		"/",
		"", // current domain
		config.GetConfig().GO_ENV == "production", // use secure in production mode
		true,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
}

func BeginOAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.Redirect(http.StatusPermanentRedirect, "/error")
		return
	}

	request := c.Request.WithContext(context.WithValue(c.Request.Context(), gothic.ProviderParamKey, provider))
	gothic.BeginAuthHandler(c.Writer, request)
}

func CompleteOAuth(c *gin.Context) {
	provider := c.Param("provider")
	request := c.Request.WithContext(context.WithValue(c.Request.Context(), gothic.ProviderParamKey, provider))

	guser, err := gothic.CompleteUserAuth(c.Writer, request)
	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "/error")
		return
	}

	user := models.UserModel{
		Email:          guser.Email,
		Name:           pickNonEmpty(guser.Name, guser.NickName, guser.Email),
		HashedPassword: "",
		OAuthUser:      true,
	}

	if !UserExists(user.Email) {
		database.DB.Create(&user)
	}

	setAuthTokenCookie(user, c)
	c.Redirect(http.StatusFound, "/")
}

func pickNonEmpty(options ...string) string {
	for _, v := range options {
		if v != "" {
			return v
		}
	}
	return ""
}
