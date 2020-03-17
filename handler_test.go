package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func (hs *HandlersSuit) SetupTest() {
	hs.db = &MapStore{tweets: make(map[string][]Tweet), userID: 0, tweetID: 0}

}

type jwtCustomInfo struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

type jwtCustomClaims struct {
	*jwt.StandardClaims
	jwtCustomInfo
}

type HandlersSuit struct {
	db *MapStore
	suite.Suite
}

func (hs *HandlersSuit) TestSignUp() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(`{    "login": "llw","password": "123","name": "Oleg","surname": "Timofeev" }`))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(hs.T(), signUp(c)) {
		assert.Equal(hs.T(), http.StatusOK, rec.Code)
		body := rec.Body.Bytes()
		var token map[string]string
		err := json.Unmarshal(body, &token)
		assert.NoError(hs.T(), err)
		validKey := []byte("secret")
		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     &jwtCustomClaims{},
			SigningKey: validKey,
		})(handler)
		makeReq := func(token string) echo.Context {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+token)
			c := e.NewContext(req, res)
			assert.NoError(hs.T(), h(c))
			return c
		}
		c := makeReq(token["token"])
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		assert.Equal(hs.T(), claims.Login, "llw")
	}
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuit))
}
