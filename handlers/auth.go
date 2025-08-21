package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"github.com/numanijaz/tinyurl/config"
	"github.com/numanijaz/tinyurl/models"
	"golang.org/x/crypto/bcrypt"
)

var userDB = make(map[string]models.UserModel)

func UserExists(username string) bool {
	if _, ok := userDB[username]; ok {
		return true
	}
	return false
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
			"sub": user.Email,
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
		"",    // current domain
		false, // TODO: change to true in production with https
		true,
	)
}

func RegisterUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	fmt.Printf("email: %s, password: %s", email, password)
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
	userDB[user.Email] = user

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

	user, userExists := userDB[email]
	if !userExists {
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
	user, exists := userDB[subject.(string)]
	if !exists {
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
		"",    // current domain
		false, // TODO: change to true in production with https
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
	}

	if !UserExists(user.Email) {
		userDB[user.Email] = user
	}
	// Save essential user fields in our own session
	// sess, _ := config.CookieStore.Get(c.Request, "app_session")
	// sess.Values["user_id"] = guser.UserID
	// sess.Values["name"] = pickNonEmpty(guser.Name, guser.NickName, guser.Email)
	// sess.Values["email"] = guser.Email
	// sess.Values["avatar"] = guser.AvatarURL
	// sess.Values["provider"] = guser.Provider
	// if err := sessions.Save(c.Request, c.Writer); err != nil {
	// 	log.Printf("session save error: %v", err)
	// }

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
