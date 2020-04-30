package main

import (
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	"time"
	"twitter/twitter/models"
	"twitter/twitter/restapi/operations/description"
)

var loginFunc = func(params description.SignInParams) middleware.Responder {
	return signIn(params)
}

func getUserTweets(params description.GetAuthorsTweetsByIDParams) middleware.Responder {
	userTweets, err := db.GetUserTweets(params.AuthorID)
	if len(userTweets) == 0 || err != nil {
		return middleware.Error(404, errNoTweet)
	}
	return description.NewGetAuthorsTweetsByIDOK().WithPayload(tweetArrayToModel(userTweets))
}

func createTweet(params description.CreateTweetParams, principal interface{}) middleware.Responder {
	user := principal.(*models.User)
	modelTweet := params.Tweet
	var tweet Tweet
	tweet.Text = modelTweet.Text
	us := new(User)
	us.ID = int(user.ID)
	us.Login = user.Login
	us = getUser(us)
	if us == nil {
		return middleware.Error(404, errNoAuth)
	}
	tweet.Time = time.Now()
	tweet.Author = us.Name + " " + us.Surname
	twt, err := db.AddTweet(&tweet, us)
	if err != nil {
		return middleware.Error(400, "Error with DB")
	}
	return description.NewCreateTweetOK().WithPayload(twt.toModel())
}

func getUser(user *User) *User {
	us, err := db.GetUserByID(strconv.Itoa(user.ID))
	if err != nil {
		return nil
	}
	return us
}
func deleteTweet(id string, user *models.User) middleware.Responder {
	us := new(User)
	us.ID = int(user.ID)
	us.Login = user.Login
	if flag, err := db.DeleteTweet(id, getUser(us)); flag == true && err == nil {
		return description.NewDeleteTweetOK()
	}

	return middleware.Error(404, errNoTweet)
}

func updateTweet(params description.UpdateTweetParams, principal interface{}) middleware.Responder {
	user := principal.(*models.User)
	us := new(User)
	us.ID = int(user.ID)
	us.Login = user.Login
	us = getUser(us)
	if flag, err := db.UpdateTweet(params.TweetID, us, params.Tweet.Text); flag == true && err == nil {
		return description.NewUpdateTweetOK()
	}
	return middleware.Error(404, errUnable)
}

func getTweet(params description.GetTweetByIDParams) middleware.Responder {
	if twt, err := db.GetTweet(params.TweetID); twt == nil || err != nil {
		return middleware.Error(404, err)
	} else {
		return description.NewGetTweetByIDOK().WithPayload(twt.toModel())
	}
}
