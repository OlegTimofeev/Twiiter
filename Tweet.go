package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Tweet struct {
	ID     string `json:"id"`
	Time   string `json:"time"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

type jwtUserClaim struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	jwt.StandardClaims
}

var tweets = make(map[string][]Tweet)
var users []User
var currentUser = new(User)
var errNoAuth = Error{Name: "Not Auth", Description: "Auth to create,update or delete twit"}
var errNoTweet = Error{Name: "No Tweet", Description: "Not find tweet "}

var mySigningKey = []byte("secret")

func GetToken(c echo.Context) error {
	claims := &jwtUserClaim{
		currentUser.ID,
		currentUser.Login,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	currentUser = new(User)
	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}
func checkUser(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtUserClaim)
	id := claims.ID
	for _, item := range users {
		if item.ID == id {
			*currentUser = item
		}
	}
	return currentUser.ID == ""
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
	if checkUser(c) {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	var tweet Tweet
	json.NewDecoder(c.Request().Body).Decode(&tweet)
	tweet.ID = strconv.Itoa(rand.Intn(10000))
	tweet.Time = time.Now().Format("2006-01-02 15:04")
	tweet.Author = currentUser.Name + " " + currentUser.Surname
	tweets[currentUser.ID] = append(tweets[currentUser.ID], tweet)
	return c.JSON(http.StatusOK, tweet)
}
func updateReplaceTweet(tweet Tweet) {
	twts := tweets[currentUser.ID]
	for index, twt := range twts {
		if twt.ID == tweet.ID {
			twts = append(twts[:index], twts[index+1:]...)
			twts = append(twts, tweet)
			tweets[currentUser.ID] = twts
		}
	}
}

func deleteTweet(c echo.Context) error {
	if checkUser(c) {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	for index, twt := range tweets[currentUser.ID] {
		if twt.ID == c.Param("id") {
			twts := append(tweets[currentUser.ID][:index], tweets[currentUser.ID][index+1:]...)
			tweets[currentUser.ID] = twts
			return c.JSON(http.StatusOK, tweets[currentUser.ID])
		}
	}
	return c.JSON(http.StatusOK, errNoTweet)
}

func updateTweet(c echo.Context) error {
	if checkUser(c) {
		return c.JSON(http.StatusOK, errNoAuth)
	}
	for _, twt := range tweets[currentUser.ID] {
		if twt.ID == c.Param("id") {
			var tweet Tweet
			json.NewDecoder(c.Request().Body).Decode(&tweet)
			tweet.ID = c.Param("id")
			tweet.Author = currentUser.Name + " " + currentUser.Surname
			tweet.Time = time.Now().Format("2006-01-02 15:04")
			updateReplaceTweet(tweet)
			return c.JSON(http.StatusOK, tweet)
		}
	}
	return c.JSON(http.StatusOK, errNoTweet)
}

func signUp(c echo.Context) error {
	currentUser.ID = strconv.Itoa(rand.Intn(1000))
	er := json.NewDecoder(c.Request().Body).Decode(&currentUser)
	if er == nil {
		users = append(users, *currentUser)
		tweets[currentUser.ID] = nil
		return GetToken(c)
	}
	return c.JSON(http.StatusOK, errNoAuth)
}

func signIn(c echo.Context) error {
	if currentUser.ID == "" {
		var loginUser User
		json.NewDecoder(c.Request().Body).Decode(&loginUser)
		for _, bdUser := range users {
			if loginUser.Login == bdUser.Login && loginUser.Password == bdUser.Password {
				return GetToken(c)
			}
		}
	}
	return c.JSON(http.StatusOK, errNoAuth)
}

func main() {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

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
