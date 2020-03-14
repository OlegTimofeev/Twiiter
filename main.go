package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"time"
)

var tweets = make(map[string][]Tweet)
var users []User
var errNoAuth = Error{Name: "Not Auth", Description: "Auth to create,update or delete twit"}
var errNoTweet = Error{Name: "No Tweet", Description: "Not find tweet "}

func initData() {
	us1 := User{ID: "3", Login: "www", Password: "123", Name: "Ol", Surname: "eg"}
	us2 := User{ID: "4", Login: "www", Password: "123", Name: "Da", Surname: "ria"}
	var twts1 []Tweet
	var twts2 []Tweet
	twts1 = append(twts1, Tweet{ID: "1", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u", Author: "Daria"})
	twts2 = append(twts2, Tweet{ID: "2", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u 2 Daria", Author: "Oleg"})
	users = append(users, us1)
	users = append(users, us2)
	tweets[us2.ID] = twts1
	tweets[us1.ID] = twts2
}

func main() {
	initData()
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	e := r.Group("/tweets")
	config := middleware.JWTConfig{
		Claims:     &jwtUserClaim{},
		SigningKey: []byte("secret"),
	}
	e.Use(middleware.JWTWithConfig(config))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("", createTweet)
	e.PUT("/:id", updateTweet)
	e.DELETE("/:id", deleteTweet)
	r.GET("/main/author/:author", getUserTweets)
	r.GET("/main", getTweets)
	r.GET("/main/:id", getTweet)
	r.POST("/signUp", signUp)
	r.POST("/signIn", signIn)

	log.Fatal(http.ListenAndServe(":8000", r))
}
