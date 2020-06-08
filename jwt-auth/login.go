package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func loginLogic(user User) (token string, err error) {

	// check if email and password match
	// ....

	// we should create proper access token
	// this access token saved as key in redis
	accessToken := "qwertyuio"

	token, err = createToken(user.Email, accessToken)
	if err != nil {
		return
	}

	// save to redis
	client := getRedisClient()
	err = client.Set(fmt.Sprintf("access_token:%s", accessToken), token, time.Duration(time.Hour)).Err()
	if err != nil {
		return
	}

	return
}
func createToken(email string, accessToken string) (token string, err error) {
	// create token claims
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["access_token"] = accessToken
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = at.SignedString([]byte("mysecretcode"))
	if err != nil {
		return
	}
	return
}
