package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"log"
	"strconv"
	"time"
	"twitter/twitter/models"
	"twitter/twitter/restapi"
	"twitter/twitter/restapi/operations"
	"twitter/twitter/restapi/operations/description"
)

var errNoAuth = Alert{Name: "Not Auth", Description: "Auth to create,update or delete twit"}
var errNoTweet = Alert{Name: "No Tweet", Description: "Not find tweet "}
var errUnable = Alert{Name: "Unable", Description: "Unable to finish operation"}
var ok = Alert{Name: "OK", Description: "Operation finished"}
var db *PostgresDB

func initData() {
	db = &PostgresDB{pgdb: nil}
	us1 := &User{Login: "log", Password: "123", Name: "Ol", Surname: "eg"}
	us2 := &User{Login: "login", Password: "123", Name: "Da", Surname: "ria"}
	tweet1 := Tweet{Time: time.Now(), Text: "text", Author: "Da ria"}
	tweet2 := Tweet{Time: time.Now(), Text: "text", Author: "Ol eg"}
	err := db.InitDB()
	panicIf(err)
	_, err = db.AddUser(us2)
	panicIf(err)
	_, err = db.AddUser(us1)
	panicIf(err)
	_, err = db.AddTweet(&tweet1, us2)
	panicIf(err)
	_, err = db.AddTweet(&tweet2, us1)
	panicIf(err)

}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func initSWHandler() *restapi.Server {
	initData()
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewTrustedTokenAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.Port = 8080

	api.APIKeyAuthAuth = func(token string) (interface{}, error) {
		claims := &jwtUserClaim{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
		if err != nil {
			return nil, err
		}
		user := new(models.User)
		id, _ := strconv.Atoi(claims.ID)
		user.ID = int64(id)
		user.Login = claims.Login

		return user, err
	}
	api.DescriptionGetAuthorsTweetsByIDHandler = description.GetAuthorsTweetsByIDHandlerFunc(func(params description.GetAuthorsTweetsByIDParams) middleware.Responder {
		return getUserTweets(params)
	})
	api.DescriptionGetTweetByIDHandler = description.GetTweetByIDHandlerFunc(func(params description.GetTweetByIDParams) middleware.Responder {
		return getTweet(params)
	})
	api.DescriptionUpdateTweetHandler = description.UpdateTweetHandlerFunc(func(params description.UpdateTweetParams, principal interface{}) middleware.Responder {
		return updateTweet(params, principal)
	})
	api.DescriptionCreateTweetHandler = description.CreateTweetHandlerFunc(func(params description.CreateTweetParams, principal interface{}) middleware.Responder {
		return createTweet(params, principal)
	})
	api.DescriptionSignUpHandler = description.SignUpHandlerFunc(func(params description.SignUpParams) middleware.Responder {

		return signUp(params)
	})

	api.DescriptionSignInHandler = description.SignInHandlerFunc(loginFunc)
	api.DescriptionDeleteTweetHandler = description.DeleteTweetHandlerFunc(func(params description.DeleteTweetParams, principal interface{}) middleware.Responder {
		user := principal.(*models.User)
		return deleteTweet(params.TweetID, user)
	})
	server.ConfigureFlags()
	server.ConfigureAPI()
	return server
}
