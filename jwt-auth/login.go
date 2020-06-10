package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type SavedToken struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func loginLogic(user User) (token string, err error) {

	// check if email and password match
	// ....

	// if match, then fetch user data from DB
	// ....

	// then we got user id
	userID := "ad82090f-0d9c-49c6-afcf-009c8df2170e"

	// we should create proper access token
	// this access token saved as key in redis
	uuid := uuid.New()
	accessToken := uuid.String()

	// create token jwt base
	token, err = createToken(user.Email, accessToken)
	if err != nil {
		return
	}

	savedToken := SavedToken{
		Token:  token,
		UserID: userID,
	}

	err = saveTokenToRedis(accessToken, savedToken)
	if err != nil {
		return
	}

	return
}
func createToken(email string, accessToken string) (token string, err error) {

	// create token claims
	// jwt extraction should have this value
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["access_token"] = accessToken
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signed string with secret code
	token, err = at.SignedString([]byte("mysecretcode"))
	if err != nil {
		return
	}
	return
}

func saveTokenToRedis(accessToken string, savedToken SavedToken) (err error) {
	client := getRedisClient()

	tokenByte, err := json.Marshal(savedToken)
	if err != nil {
		return
	}

	// we save this value in redis
	// access_token:YOUR_ACCESS_TOKEN : {
	// Token : qwertyui,
	// UserID : 1234567
	// }

	err = client.Set(fmt.Sprintf("access_token:%s", accessToken), tokenByte, time.Duration(time.Hour)).Err()
	if err != nil {
		return
	}

	return
}
