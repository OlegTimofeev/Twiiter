package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var mySigningKey = []byte("secret")

func GetToken(c echo.Context, us User) error {
	claims := &jwtUserClaim{
		us.ID,
		us.Login,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}
func signUp(c echo.Context) error {
	us := new(User)
	us.ID = strconv.Itoa(rand.Intn(1000))
	er := json.NewDecoder(c.Request().Body).Decode(&us)
	if er == nil {
		users = append(users, *us)
		tweets[us.ID] = nil
		return GetToken(c, *us)
	}
	return c.JSON(http.StatusOK, errNoAuth)
}

func signIn(c echo.Context) error {
	var loginUser User
	json.NewDecoder(c.Request().Body).Decode(&loginUser)
	for _, bdUser := range users {
		if loginUser.Login == bdUser.Login && loginUser.Password == bdUser.Password {
			return GetToken(c, bdUser)
		}
	}
	return c.JSON(http.StatusOK, errNoAuth)
}
