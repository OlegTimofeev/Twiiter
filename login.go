package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	"time"
	"twitter/twitter/models"
)

var mySigningKey = []byte("secret")

func userTokenResponse(us *User) middleware.Responder {
	claims := &jwtUserClaim{
		ID:    strconv.Itoa(us.ID),
		Login: us.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	tkn := models.Token{
		Token: tokenString,
	}
	return middleware.Error(200, tkn)
}

func signUp(user *User) middleware.Responder {
	user, err := db.AddUser(user)
	if err != nil {
		return userTokenResponse(nil)
	}
	return userTokenResponse(user)
}

func signIn(user *User) middleware.Responder {
	loginUser := user
	if us, isFinded, err := db.CheckLoginPassword(loginUser.Login, loginUser.Password); isFinded && err == nil {
		return userTokenResponse(us)
	}
	return middleware.Error(404, nil)
}
