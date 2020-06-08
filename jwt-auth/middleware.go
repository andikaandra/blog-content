package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenClaims struct {
	Email       string
	AccessToken string
	Exp         int64
}

func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValid := isTokenValid(c)
		if !tokenValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"data":   "unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func isTokenValid(c *gin.Context) bool {
	tokenString, err := getToken(c)
	if err != nil {
		return false
	}

	token, err := validateToken(tokenString)
	if err != nil {
		return false
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false
	}

	return true
}

func getToken(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("user_session")
	if err != nil {
		return "", err
	}

	return cookie, nil
}

func validateToken(tokenString string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("mysecretcode"), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func getTokenData(c *gin.Context) (tokenClaims TokenClaims, err error) {
	tokenString, err := getToken(c)
	if err != nil {
		return
	}

	token, err := validateToken(tokenString)
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return
		}

		accessToken, ok := claims["acces_token"].(string)
		if !ok {
			return
		}

		exp, ok := claims["exp"].(int64)
		if !ok {
			return
		}

		tokenClaims = TokenClaims{
			Email:       email,
			AccessToken: accessToken,
			Exp:         exp,
		}
		return
	}

	return
}
