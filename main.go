package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"twitter/twitter/restapi"
	"twitter/twitter/restapi/operations"
	"twitter/twitter/restapi/operations/description"
)

var errNoAuth = Alert{Name: "Not Auth", Description: "Auth to create,update or delete twit"}
var errNoTweet = Alert{Name: "No Tweet", Description: "Not find tweet "}
var errUnable = Alert{Name: "Unable", Description: "Unable to finish operation"}
var errBadReq = Alert{Name: "Bad Request", Description: "Error.Bad Request"}
var ok = Alert{Name: "OK", Description: "Operation finished"}
var db *PostgresDB

func initData() {
	db = &PostgresDB{pgdb: nil}
	/*db = &MapStore{tweets: make(map[string][]*Tweet), userID: 1, tweetID: 1}*/
	us1 := &User{Login: "www", Password: "123", Name: "Ol", Surname: "eg"}
	us2 := &User{Login: "wwww", Password: "123", Name: "Da", Surname: "ria"}
	tweet1 := Tweet{Time: time.Now(), Text: "I love u", Author: "Daria"}
	tweet2 := Tweet{Time: time.Now(), Text: "I love u 2 Daria", Author: "Oleg"}
	//db.AddTweet(&tweet1, db.AddUser(&us2))
	//db.AddTweet(&tweet2, db.AddUser(&us1))
	db.InitDB()
	_, err := db.AddUser(us2)
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

func deleteTweet(c echo.Context) error {
	if flag, err := db.DeleteTweet(c.Param("id"), getUser(c)); flag == true && err == nil {
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
	if flag, err := db.UpdateTweet(c.Param("id"), us, twt.Text); flag == true && err == nil {
		return c.JSON(http.StatusOK, ok)
	}
	return c.JSON(http.StatusOK, errUnable)
}

func getUser(c echo.Context) *User {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtUserClaim)
	id := claims.ID
	us, err := db.GetUserByID(id)
	if err != nil {
		return nil
	}
	return us
}

func getTweets(c echo.Context) error {
	tweets, err := db.GetTweets()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errNoTweet)
	}
	return c.JSON(http.StatusOK, tweets)
}

func getTweet(c echo.Context) error {
	if twt, err := db.GetTweet(c.Param("id")); twt == nil || err != nil {
		return c.JSON(http.StatusOK, errNoTweet)
	} else {
		return c.JSON(http.StatusOK, twt)
	}
}

func getUserTweets(c echo.Context) error {
	userTweets, err := db.GetUserTweets(c.Param("authorID"))
	if len(userTweets) == 0 || err != nil {
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
	twt, err := db.AddTweet(&tweet, us)
	if err != nil {
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusOK, twt)
}

func initHandler() http.Handler {
	initData()
	r := echo.New()
	//r.Use(middleware.Logger())
	//r.Use(middleware.Recover())
	e := r.Group("/tweets")
	//config := middleware.JWTConfig{
	//	Claims:     &jwtUserClaim{},
	//	SigningKey: []byte("secret"),
	//}
	//e.Use(middleware.JWTWithConfig(config))
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
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
	//r := initHandler()
	//log.Fatal(http.ListenAndServe(":8000", r))
	initData()
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewTrustedTokenAPI(swaggerSpec)
	server := restapi.NewServer(api)

	api.DescriptionGetTweetByIDHandler = description.GetTweetByIDHandlerFunc(func(params description.GetTweetByIDParams) middleware.Responder {
		if twt, err := db.GetTweet(params.TweetID); twt == nil || err != nil {
			return middleware.Error(404, err)
		} else {
			return middleware.Error(200, twt)
		}
	})
	api.DescriptionCreateTweetHandler = description.CreateTweetHandlerFunc(func(params description.CreateTweetParams) middleware.Responder {
		return middleware.NotImplemented("not implemented")
	})
	api.DescriptionSignUpHandler = description.SignUpHandlerFunc(func(params description.SignUpParams) middleware.Responder {
		us := params.User
		claims := &jwtUserClaim{
			ID:    strconv.Itoa(int(us.ID)),
			Login: us.Login,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(mySigningKey)
		return middleware.Error(200, tokenString)
	})
	api.DescriptionSignInHandler = description.SignInHandlerFunc(func(params description.SignInParams) middleware.Responder {
		return middleware.NotImplemented("not implemented")
	})

	defer server.Shutdown()
	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Trusted Token API"
	parser.LongDescription = "This is a license API in cloud for AxxonNext"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
