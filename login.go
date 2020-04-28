package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	"time"
	"twitter/twitter/models"
	"twitter/twitter/restapi/operations/description"
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

func signUp(params description.SignUpParams) middleware.Responder {
	newUser := new(User)
	newUser.Login = params.User.Login
	newUser.Password = params.User.Password
	newUser.Name = params.User.Name
	newUser.Surname = params.User.Surname
	newUser, err := db.AddUser(newUser)
	if err != nil {
		return userTokenResponse(nil)
	}
	return userTokenResponse(newUser)
}

func signIn(params description.SignInParams) middleware.Responder {
	user := params.User
	var loginUser User
	loginUser.Login = user.Login
	loginUser.Password = user.Password
	if us, isFinded, err := db.CheckLoginPassword(loginUser.Login, loginUser.Password); isFinded && err == nil {
		return userTokenResponse(us)
	}
	return middleware.Error(404, nil)
}
