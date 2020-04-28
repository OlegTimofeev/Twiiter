package main

import (
	client2 "github.com/go-openapi/runtime/client"
	util2 "github.com/itimofeev/go-util"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strconv"
	"testing"
	"twitter/twitter/client"
	"twitter/twitter/client/description"
	"twitter/twitter/models"
)

func (hs *HandlersSuit) SetupTest() {
	httpClient := &http.Client{Transport: util2.NewTransport(initSWHandler().GetHandler())}
	c := client2.NewWithClient(client.DefaultHost, client.DefaultBasePath, client.DefaultSchemes, httpClient)
	hs.deviceRegistry = client.New(c, nil)
	//from initialization
	hs.password = "123"
	hs.login = "login"
}

type HandlersSuit struct {
	login          string
	password       string
	text           string
	deviceRegistry *client.TrustedToken
	suite.Suite
}

func (hs *HandlersSuit) TestSignUpAndCreateTweet() {
	signupOK, err := hs.deviceRegistry.Description.SignUp(description.NewSignUpParams().WithUser(description.SignUpBody{
		Login:    "llw",
		Password: "123",
		Name:     "Ole",
		Surname:  "G",
	}))
	hs.Require().NoError(err)
	token := signupOK.Payload.Token
	createTweetOk, err := hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "12333333333",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
}

func (hs *HandlersSuit) TestSignInAndCreateTweet() {
	signinOK, err := hs.deviceRegistry.Description.SignIn(description.NewSignInParams().WithUser(description.SignInBody{
		Login:    hs.login,
		Password: hs.password,
	}))
	hs.Require().NoError(err)
	token := signinOK.Payload.Token
	createTweetOk, err := hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "12333333333",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
}

func (hs *HandlersSuit) TestGetUserTweets() {
	signupOK, err := hs.deviceRegistry.Description.SignUp(description.NewSignUpParams().WithUser(description.SignUpBody{
		Login:    "llw",
		Password: "123",
		Name:     "Ole",
		Surname:  "G",
	}))
	hs.Require().NoError(err)
	token := signupOK.Payload.Token
	createTweetOk, err := hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "first tweet",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
	createTweetOk, err = hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "second tweet",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
	userTweets, err := hs.deviceRegistry.Description.GetAuthorsTweetsByID(description.NewGetAuthorsTweetsByIDParams().WithAuthorID(strconv.Itoa(int(createTweetOk.Payload.AuthorID))))
	hs.Require().NoError(err)
	countOfCreatedTweets := 2
	hs.Require().NotNil(userTweets)
	hs.Require().Equal(countOfCreatedTweets, len(userTweets.Payload))
}

func (hs *HandlersSuit) TestDeleteTweet() {
	signupOK, err := hs.deviceRegistry.Description.SignUp(description.NewSignUpParams().WithUser(description.SignUpBody{
		Login:    "llw",
		Password: "123",
		Name:     "Ole",
		Surname:  "G",
	}))
	hs.Require().NoError(err)
	token := signupOK.Payload.Token
	createTweetOk, err := hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "first tweet",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
	userTweets, err := hs.deviceRegistry.Description.GetAuthorsTweetsByID(description.NewGetAuthorsTweetsByIDParams().WithAuthorID(strconv.Itoa(int(createTweetOk.Payload.AuthorID))))
	hs.Require().NoError(err)
	countOfCreatedTweets := 1
	hs.Require().NotNil(userTweets)
	hs.Require().Equal(countOfCreatedTweets, len(userTweets.Payload))
	_, err = hs.deviceRegistry.Description.DeleteTweet(description.NewDeleteTweetParams().WithTweetID(strconv.Itoa(int(createTweetOk.Payload.ID))), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
}

func (hs *HandlersSuit) TestDeleteTweetError() {
	signupOK, err := hs.deviceRegistry.Description.SignUp(description.NewSignUpParams().WithUser(description.SignUpBody{
		Login:    "llw",
		Password: "123",
		Name:     "Ole",
		Surname:  "G",
	}))
	hs.Require().NoError(err)
	token := signupOK.Payload.Token
	createTweetOk, err := hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "first tweet",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
	userTweets, err := hs.deviceRegistry.Description.GetAuthorsTweetsByID(description.NewGetAuthorsTweetsByIDParams().WithAuthorID(strconv.Itoa(int(createTweetOk.Payload.AuthorID))))
	hs.Require().NoError(err)
	countOfCreatedTweets := 1
	hs.Require().NotNil(userTweets)
	hs.Require().Equal(countOfCreatedTweets, len(userTweets.Payload))
	signinOK, err := hs.deviceRegistry.Description.SignIn(description.NewSignInParams().WithUser(description.SignInBody{
		Login:    hs.login,
		Password: hs.password,
	}))
	hs.Require().NoError(err)
	token = signinOK.Payload.Token
	_, err = hs.deviceRegistry.Description.DeleteTweet(description.NewDeleteTweetParams().WithTweetID(strconv.Itoa(int(createTweetOk.Payload.ID))), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().Error(err)

}

func (hs *HandlersSuit) TestUpdateTweet() {
	signupOK, err := hs.deviceRegistry.Description.SignUp(description.NewSignUpParams().WithUser(description.SignUpBody{
		Login:    "llw",
		Password: "123",
		Name:     "Ole",
		Surname:  "G",
	}))
	hs.Require().NoError(err)
	token := signupOK.Payload.Token
	createTweetOk, err := hs.deviceRegistry.Description.CreateTweet(description.NewCreateTweetParams().WithTweet(&models.Tweet{
		Text: "first tweet",
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.Require().NoError(err)
	hs.Require().NotNil(createTweetOk)
	userTweets, err := hs.deviceRegistry.Description.GetAuthorsTweetsByID(description.NewGetAuthorsTweetsByIDParams().WithAuthorID(strconv.Itoa(int(createTweetOk.Payload.AuthorID))))
	hs.Require().NoError(err)
	countOfCreatedTweets := 1
	hs.Require().NotNil(userTweets)
	hs.Require().Equal(countOfCreatedTweets, len(userTweets.Payload))
	_, err = hs.deviceRegistry.Description.UpdateTweet(description.NewUpdateTweetParams().WithTweetID(strconv.Itoa(int(createTweetOk.Payload.ID))).WithTweet(&models.Tweet{
		Text: hs.text,
	}), client2.APIKeyAuth("Authorization", "header", token))
	hs.NoError(err)
	updatedTweet, err := hs.deviceRegistry.Description.GetTweetByID(description.NewGetTweetByIDParams().WithTweetID(strconv.Itoa(int(createTweetOk.Payload.ID))))
	hs.Require().NoError(err)
	hs.Require().Equal(hs.text, updatedTweet.Payload.Text)

}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuit))
}
