package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()

	r.POST("/login", loginRoute)
	r.GET("/users", allUserRoute)

	r.Use(tokenAuthMiddleware())
	{
		r.GET("/books", allBookRoute)
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

	// successfull login, save toke as cookies
	c.SetCookie("user_session", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"data": token,
	})
}

func allUserRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "all users, everybody can access",
	})
}

func allBookRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "all books, only authenticated user can access",
	})
}

func findUserRoute(c *gin.Context) {
	userID := c.Param("userID")

	tokenClaims, err := getTokenData(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Println(tokenClaims)

	allowed := isResourceOwner(tokenClaims, userID)
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you're not the resource owner",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fmt.Sprintf("data user : %s, only the user can access", userID),
	})
}
