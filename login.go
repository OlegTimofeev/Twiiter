package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

var mySigningKey = []byte("secret")

func userTokenResponse(c echo.Context, us User) error {
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
	mutex.Lock()
	us.ID = strconv.Itoa(getUserId())
	mutex.Unlock()
	er := json.NewDecoder(c.Request().Body).Decode(&us)
	if er == nil {
		db.addUser(*us)
		return userTokenResponse(c, *us)
	}
	return c.JSON(http.StatusOK, errNoAuth)
}

func signIn(c echo.Context) error {
	var loginUser User
	json.NewDecoder(c.Request().Body).Decode(&loginUser)
	if us, isFinded := db.checkLoginPassword(loginUser.Login, loginUser.Password); isFinded {
		return userTokenResponse(c, *us)
	}
	return c.JSON(http.StatusOK, errNoAuth)
}
