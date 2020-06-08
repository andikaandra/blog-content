package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()

	r.POST("/login", loginRoute)
	r.GET("/users", allUserRoute)

	r.Use(tokenAuthMiddleware())
	{
		r.GET("/user/:userID", findUserRoute)
	}

	r.Run()
}

func loginRoute(c *gin.Context) {
	var user User

	c.Bind(&user)
	token, err := loginLogic(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	c.SetCookie("user_session", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"data": token,
	})
}

func allUserRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "all user",
	})
}

func findUserRoute(c *gin.Context) {
	userID := c.Param("userID")

	c.JSON(http.StatusOK, gin.H{
		"data": "data user :" + userID,
	})
}
