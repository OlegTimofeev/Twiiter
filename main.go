package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var tweetId = 3
var userId = 5
var mutex = &sync.Mutex{}
var errNoAuth = Alert{Name: "Not Auth", Description: "Auth to create,update or delete twit"}
var errNoTweet = Alert{Name: "No Tweet", Description: "Not find tweet "}
var errUndable = Alert{Name: "Unable", Description: "Unable to finish operation"}
var ok = Alert{Name: "OK", Description: "Operation finished"}

var db = MapStore{tweets: make(map[string][]Tweet)}

func initData() {
	us1 := User{ID: "3", Login: "www", Password: "123", Name: "Ol", Surname: "eg"}
	us2 := User{ID: "4", Login: "wwww", Password: "123", Name: "Da", Surname: "ria"}
	db.addUser(us1)
	db.addUser(us2)
	db.addTweet(Tweet{ID: "1", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u", Author: "Daria"}, us2)
	db.addTweet(Tweet{ID: "2", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u 2 Daria", Author: "Oleg"}, us1)
}

func getTweetId() int {
	returnId := tweetId
	tweetId++
	return returnId
}

func getUserId() int {
	returnId := userId
	userId++
	return returnId
}

func deleteTweet(c echo.Context) error {
	if db.deleteTweet(c.Param("id"), getUser(c).ID) {
		return c.JSON(http.StatusOK, ok)
	}

	return c.JSON(http.StatusOK, errNoTweet)
}

func updateTweet(c echo.Context) error {
	us := getUser(c)
	var twt Tweet
	json.NewDecoder(c.Request().Body).Decode(&twt)
	if db.updateTweet(c.Param("id"), us.ID, twt.Text) {
		return c.JSON(http.StatusOK, ok)
	}
	return c.JSON(http.StatusOK, errUndable)
}

func getUser(c echo.Context) *User {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtUserClaim)
	id := claims.ID
	return db.getUserById(id)
}

func getTweets(c echo.Context) error {
	c.JSON(http.StatusOK, db.getTweets())
	return c.String(http.StatusOK, "That's all folks")
}

func getTweet(c echo.Context) error {
	if twt := db.getTweet(c.Param("id")); twt == nil {
		return c.JSON(http.StatusOK, errNoTweet)
	} else {
		return c.JSON(http.StatusOK, twt)
	}
}

func getUserTweets(c echo.Context) error {
	userTweets := db.getUserTweets(c.Param("authorId"))
	if len(userTweets) == 0 {
		return c.JSON(http.StatusOK, errNoTweet)
	}
	return c.JSON(http.StatusOK, userTweets)
}

func createTweet(c echo.Context) error {
	us := getUser(c)
	if us == nil {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	var tweet Tweet
	json.NewDecoder(c.Request().Body).Decode(&tweet)
	mutex.Lock()
	tweet.ID = strconv.Itoa(getTweetId())
	mutex.Unlock()
	tweet.Time = time.Now().Format("2006-01-02 15:04")
	tweet.Author = us.Name + " " + us.Surname
	tweet.AuthorID = us.ID
	db.addTweet(tweet, *us)
	return c.JSON(http.StatusOK, tweet)
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
	r.GET("/main/author/:authorId", getUserTweets)
	r.GET("/main", getTweets)
	r.GET("/main/:id", getTweet)
	r.POST("/signUp", signUp)
	r.POST("/signIn", signIn)

	log.Fatal(http.ListenAndServe(":8000", r))
}
