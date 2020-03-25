package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

var errNoAuth = Alert{Name: "Not Auth", Description: "Auth to create,update or delete twit"}
var errNoTweet = Alert{Name: "No Tweet", Description: "Not find tweet "}
var errUnable = Alert{Name: "Unable", Description: "Unable to finish operation"}
var errBadReq = Alert{Name: "Bad Request", Description: "Error.Bad Request"}
var ok = Alert{Name: "OK", Description: "Operation finished"}
var db *MapStore

func initData() {
	db = &MapStore{tweets: make(map[string][]*Tweet), userID: 1, tweetID: 1}
	us1 := User{Login: "www", Password: "123", Name: "Ol", Surname: "eg"}
	us2 := User{Login: "wwww", Password: "123", Name: "Da", Surname: "ria"}
	db.addTweet(&Tweet{Time: time.Now(), Text: "I love u", Author: "Daria"}, db.addUser(&us2))
	db.addTweet(&Tweet{Time: time.Now(), Text: "I love u 2 Daria", Author: "Oleg"}, db.addUser(&us1))
}

func deleteTweet(c echo.Context) error {
	if db.deleteTweet(c.Param("id"), getUser(c)) {
		return c.JSON(http.StatusOK, ok)
	}

	return c.JSON(http.StatusOK, errNoTweet)
}

func updateTweet(c echo.Context) error {
	us := getUser(c)
	var twt Tweet
	err := json.NewDecoder(c.Request().Body).Decode(&twt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errBadReq)
	}
	if db.updateTweet(c.Param("id"), us, twt.Text) {
		return c.JSON(http.StatusOK, ok)
	}
	return c.JSON(http.StatusOK, errUnable)
}

func getUser(c echo.Context) *User {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtUserClaim)
	id := claims.ID
	return db.getUserByID(id)
}

func getTweets(c echo.Context) error {
	return c.JSON(http.StatusOK, db.getTweets())
}

func getTweet(c echo.Context) error {
	if twt := db.getTweet(c.Param("id")); twt == nil {
		return c.JSON(http.StatusOK, errNoTweet)
	} else {
		return c.JSON(http.StatusOK, twt)
	}
}

func getUserTweets(c echo.Context) error {
	userTweets := db.getUserTweets(c.Param("authorID"))
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
	if err := json.NewDecoder(c.Request().Body).Decode(&tweet); err != nil {
		return c.JSON(http.StatusBadRequest, errBadReq)
	}
	tweet.Time = time.Now()
	tweet.Author = us.Name + " " + us.Surname
	return c.JSON(http.StatusOK, *db.addTweet(&tweet, us))
}

func initHandler() http.Handler {
	fmt.Println("Successfully connected!")
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
	r.GET("/main/author/:authorID", getUserTweets)
	r.GET("/main", getTweets)
	r.GET("/main/:id", getTweet)
	r.POST("/signUp", signUp)
	r.POST("/signIn", signIn)
	return r
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	r := initHandler()
	log.Fatal(http.ListenAndServe(":8000", r))
}
