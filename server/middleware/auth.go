package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(e *gin.Engine) {

	router := e.Group("/auth")

	router.POST("/login", handleLogin())
	router.GET("/logout", handleLogout())
}

func handleLogin() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("id", 152456789)
		session.Set("email", "test@gmail.com")
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"message": "User signed in",
		})
	}

	return gin.HandlerFunc(fn)
}

func handleLogout() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"message": "User signed out",
		})
	}

	return gin.HandlerFunc(fn)
}

func HandleAuthVerify() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		sessionId := session.Get("id")

		// you don't need a session to log in
		if sessionId == nil && c.Request.URL.EscapedPath() != "/auth/login" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
		}
	}

	return gin.HandlerFunc(fn)
}
