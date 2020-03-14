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

func deleteTweet(c echo.Context) error {
	us := getUser(c)
	if us == nil {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	for index, twt := range tweets[us.ID] {
		if twt.ID == c.Param("id") {
			twts := append(tweets[us.ID][:index], tweets[us.ID][index+1:]...)
			tweets[us.ID] = twts
			return c.JSON(http.StatusOK, tweets[us.ID])
		}
	}
	return c.JSON(http.StatusOK, errNoTweet)
}

func updateTweet(c echo.Context) error {
	us := getUser(c)
	if us == nil {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	for _, twt := range tweets[us.ID] {
		if twt.ID == c.Param("id") {
			var tweet Tweet
			json.NewDecoder(c.Request().Body).Decode(&tweet)
			tweet.ID = c.Param("id")
			tweet.Author = us.Name + " " + us.Surname
			tweet.Time = time.Now().Format("2006-01-02 15:04")
			updateReplaceTweet(tweet, us)
			return c.JSON(http.StatusOK, tweet)
		}
	}
	return c.JSON(http.StatusOK, errNoTweet)
}

func getUser(c echo.Context) *User {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtUserClaim)
	id := claims.ID
	for _, item := range users {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func getTweets(c echo.Context) error {
	for _, us := range users {
		usersTweets := tweets[us.ID]
		c.JSON(http.StatusOK, usersTweets)
	}
	return c.String(http.StatusOK, "That's all folks")
}

func getTweet(c echo.Context) error {
	for us := range tweets {
		for _, twt := range tweets[us] {
			if twt.ID == c.Param("id") {
				return c.JSON(http.StatusOK, twt)
			}
		}
	}
	return c.JSON(http.StatusOK, errNoTweet)
}

func getUserTweets(c echo.Context) error {
	for _, us := range users {
		if us.Name+" "+us.Surname == c.Param("author") {
			for _, twt := range tweets[us.ID] {
				c.JSON(http.StatusOK, twt)
			}
			return nil
		}
	}
	return c.JSON(http.StatusOK, errNoTweet)
}

func createTweet(c echo.Context) error {
	us := getUser(c)
	if us == nil {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	var tweet Tweet
	json.NewDecoder(c.Request().Body).Decode(&tweet)
	tweet.ID = strconv.Itoa(rand.Intn(10000))
	tweet.Time = time.Now().Format("2006-01-02 15:04")
	tweet.Author = us.Name + " " + us.Surname
	tweets[us.ID] = append(tweets[us.ID], tweet)
	return c.JSON(http.StatusOK, tweet)
}
func updateReplaceTweet(tweet Tweet, us *User) {
	twts := tweets[us.ID]
	for index, twt := range twts {
		if twt.ID == tweet.ID {
			twts = append(twts[:index], twts[index+1:]...)
			twts = append(twts, tweet)
			tweets[us.ID] = twts
		}
	}
}
