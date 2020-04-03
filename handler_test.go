package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func (hs *HandlersSuit) SetupTest() {
	hs.testReaderSignUpUser = strings.NewReader(`{    "login": "llw","password": "123","name": "Oleg","surname": "Timofeev" }`)
	hs.testReaderSignInUser = strings.NewReader(`{    "login": "www","password": "123" }`)
	hs.testReaderTweet = strings.NewReader(`{"text":"text"}`)
	hs.testText = "text"
}

type HandlersSuit struct {
	testText             string
	testReaderSignInUser io.Reader
	testReaderSignUpUser io.Reader
	testReaderTweet      io.Reader
	suite.Suite
}

func userTokenResponseTest(r http.Handler, reader io.Reader, inOrUp string) (error, *string) {
	var req *http.Request
	if inOrUp == "in" {
		req = httptest.NewRequest(http.MethodPost, "/signIn", reader)
	} else {
		req = httptest.NewRequest(http.MethodPost, "/signUp", reader)
	}
	recSignUp := httptest.NewRecorder()
	r.ServeHTTP(recSignUp, req)
	body := recSignUp.Body.Bytes()
	var token Tok
	err := json.Unmarshal(body, &token)
	if err != nil {
		return err, nil
	}
	return nil, &token.TokenValue
}

func (hs *HandlersSuit) TestUserTokenResponse() {
	r := initHandler()
	err, token := userTokenResponseTest(r, hs.testReaderSignUpUser, "up")
	hs.Require().NoError(err)
	hs.Require().NotNil(token)
	err, token = userTokenResponseTest(r, hs.testReaderSignInUser, "in")
	hs.Require().NoError(err)
	hs.Require().NotNil(token)
}

func (hs *HandlersSuit) TestSignUpAndCreateTweet() {
	r := initHandler()
	err, token := userTokenResponseTest(r, hs.testReaderSignUpUser, "up")
	hs.Require().NoError(err)
	recCreate := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/tweets", hs.testReaderTweet)
	req2.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*token)
	r.ServeHTTP(recCreate, req2)
	result := recCreate.Body.Bytes()
	var tweet Tweet
	err = json.Unmarshal(result, &tweet)
	hs.Require().NoError(err)
	hs.Require().Equal(hs.testText, tweet.Text)
}

func (hs *HandlersSuit) TestSignInAndCreateTweet() {
	r := initHandler()
	err, token := userTokenResponseTest(r, hs.testReaderSignInUser, "in")
	hs.Require().NoError(err)
	recCreate := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/tweets", hs.testReaderTweet)
	req2.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*token)
	r.ServeHTTP(recCreate, req2)
	result := recCreate.Body.Bytes()
	var tweet Tweet
	err = json.Unmarshal(result, &tweet)
	hs.Require().NoError(err)
	hs.Require().Equal(hs.testText, tweet.Text)
}

func (hs *HandlersSuit) TestGetAllTweets() {
	r := initHandler()
	req := httptest.NewRequest(http.MethodGet, "/main", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	body := rec.Body.Bytes()
	var allTweets []*Tweet
	err := json.Unmarshal(body, &allTweets)
	hs.Require().NoError(err)
	hs.Require().NotEmpty(allTweets)
	hs.Require().Equal(2, len(allTweets))
}

func (hs *HandlersSuit) TestGetUserTweets() {
	r := initHandler()
	err, token := userTokenResponseTest(r, hs.testReaderSignInUser, "in")
	hs.Require().NoError(err)
	recCreate := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/tweets", hs.testReaderTweet)
	req2.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*token)
	r.ServeHTTP(recCreate, req2)
	resultCreate := recCreate.Body.Bytes()
	var tweet Tweet
	err = json.Unmarshal(resultCreate, &tweet)
	hs.Require().NoError(err)
	hs.Require().Equal(hs.testText, tweet.Text)
}

func (hs *HandlersSuit) TestDeleteTweet() {
	r := initHandler()
	errToken1, tokenUser1 := userTokenResponseTest(r, hs.testReaderSignUpUser, "up")
	hs.Require().NoError(errToken1)
	recCreate := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/tweets", hs.testReaderTweet)
	req2.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser1)
	r.ServeHTTP(recCreate, req2)
	resultCreate := recCreate.Body.Bytes()
	var tweet Tweet
	err := json.Unmarshal(resultCreate, &tweet)
	hs.Require().NoError(err)
	hs.Require().Equal(hs.testText, tweet.Text)
	recDelete := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodDelete, "/tweets/"+strconv.Itoa(tweet.ID), nil)
	req3.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser1)
	r.ServeHTTP(recDelete, req3)
	resultDelete := recDelete.Body.Bytes()
	var deleted Alert
	err = json.Unmarshal(resultDelete, &deleted)
	hs.Require().NoError(err)
	hs.Require().Equal(ok.Name, deleted.Name)
	hs.Require().Equal(ok.Description, deleted.Description)
}

func (hs *HandlersSuit) TestDeleteTweetError() {
	r := initHandler()
	errToken1, tokenUser1 := userTokenResponseTest(r, hs.testReaderSignUpUser, "up")
	errToken2, tokenUser2 := userTokenResponseTest(r, hs.testReaderSignInUser, "in")
	hs.Require().NoError(errToken1)
	hs.Require().NoError(errToken2)
	recCreate := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/tweets", hs.testReaderTweet)
	req2.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser1)
	r.ServeHTTP(recCreate, req2)
	resultCreate := recCreate.Body.Bytes()
	var tweet Tweet
	err := json.Unmarshal(resultCreate, &tweet)
	hs.Require().NoError(err)
	hs.Require().Equal(hs.testText, tweet.Text)
	recDelete := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodDelete, "/tweets/"+strconv.Itoa(tweet.ID), nil)
	req3.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser2)
	r.ServeHTTP(recDelete, req3)
	resultDelete := recDelete.Body.Bytes()
	var deleted Alert
	//пользователь который не создавал твит получит ошибку errNoTweet
	err = json.Unmarshal(resultDelete, &deleted)
	hs.Require().NoError(err)
	hs.Require().Equal(errNoTweet.Name, deleted.Name)
	hs.Require().Equal(errNoTweet.Description, deleted.Description)
	recDelete = httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodDelete, "/tweets/"+strconv.Itoa(tweet.ID), nil)
	req4.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser1)
	r.ServeHTTP(recDelete, req4)
	resultDelete = recDelete.Body.Bytes()
	//пользователь который создал твит удалит его
	err = json.Unmarshal(resultDelete, &deleted)
	hs.Require().NoError(err)
	hs.Require().Equal(ok.Name, deleted.Name)
	hs.Require().Equal(ok.Description, deleted.Description)
}

func (hs *HandlersSuit) TestUpdateTweet() {
	r := initHandler()
	errToken1, tokenUser1 := userTokenResponseTest(r, hs.testReaderSignUpUser, "up")
	hs.Require().NoError(errToken1)
	recCreate := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/tweets", hs.testReaderTweet)
	req2.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser1)
	r.ServeHTTP(recCreate, req2)
	resultCreate := recCreate.Body.Bytes()
	var tweet Tweet
	err := json.Unmarshal(resultCreate, &tweet)
	hs.Require().NoError(err)
	hs.Require().Equal(hs.testText, tweet.Text)
	recUpdate := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPut, "/tweets/"+strconv.Itoa(tweet.ID), strings.NewReader(`{"text":"update"}`))
	req3.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+*tokenUser1)
	r.ServeHTTP(recUpdate, req3)
	resultUpdate := recUpdate.Body.Bytes()
	var updateAlert Alert
	errUpdate := json.Unmarshal(resultUpdate, &updateAlert)
	hs.Require().NoError(errUpdate)
	hs.Require().Equal(ok.Name, updateAlert.Name)
	hs.Require().Equal(ok.Description, updateAlert.Description)
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuit))
}
