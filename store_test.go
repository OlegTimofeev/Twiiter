package main

import (
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

func (ts *MapStoreTestSuit) SetupTest() {
	ts.db = &MapStore{tweets: make(map[int][]*Tweet), userID: 0, tweetID: 0}
	ts.user1 = User{Login: "www", Password: "123", Name: "Ol", Surname: "eg"}
	ts.user2 = User{Login: "ufo", Password: "321", Name: "Il", Surname: "ya"}
	ts.tweet1 = Tweet{Text: "Test"}
	ts.tweet2 = Tweet{Text: "Test2"}
	ts.tweet3 = Tweet{Text: "Test3"}
	ts.newText = "New Text"
}

type MapStoreTestSuit struct {
	db *MapStore
	suite.Suite
	user1   User
	user2   User
	tweet1  Tweet
	tweet2  Tweet
	tweet3  Tweet
	newText string
}

func (ts *MapStoreTestSuit) TestAddUser() {

	createdUser := ts.db.AddUser(&ts.user1)
	us, isExist := ts.db.CheckLoginPassword(createdUser.Login, createdUser.Password)
	ts.Require().True(isExist)
	ts.Require().Equal(us.Login, createdUser.Login)
	ts.Require().Equal(us.Password, createdUser.Password)
	ts.Require().Equal(us.Name, createdUser.Name)
	ts.Require().Equal(us.Surname, createdUser.Surname)
}

func (ts *MapStoreTestSuit) TestAddTweet() {
	createdUser := ts.db.AddUser(&ts.user1)
	tweet := ts.db.AddTweet(&ts.tweet1, createdUser)
	ts.Require().Equal(*tweet, *ts.db.GetTweet(strconv.Itoa(tweet.ID)))
}

func (ts *MapStoreTestSuit) TestGetUserTweets() {
	createdUser1 := ts.db.AddUser(&ts.user1)
	ts.Require().Nil(ts.db.GetUserTweets(strconv.Itoa(createdUser1.ID)))
	createdUser2 := ts.db.AddUser(&ts.user2)
	tweet1 := ts.db.AddTweet(&ts.tweet1, createdUser1)
	tweets := ts.db.GetUserTweets(strconv.Itoa(createdUser1.ID))
	ts.Require().Contains(tweets, tweet1)
	tweet2 := ts.db.AddTweet(&ts.tweet1, createdUser1)
	tweet3 := ts.db.AddTweet(&ts.tweet3, createdUser2)
	tweets = ts.db.GetUserTweets(strconv.Itoa(createdUser1.ID))
	countOfTweets := 2
	ts.Require().NotContains(tweets, *tweet3)
	ts.Require().Contains(tweets, tweet1)
	ts.Require().Contains(tweets, tweet2)
	ts.Require().Equal(countOfTweets, len(tweets))
}

func (ts *MapStoreTestSuit) TestUpdateTweet() {
	createdUser1 := ts.db.AddUser(&ts.user1)
	tweet1 := ts.db.AddTweet(&ts.tweet1, createdUser1)
	isChanged := ts.db.UpdateTweet(strconv.Itoa(tweet1.ID), createdUser1, ts.newText)
	ts.Require().True(isChanged)
	ts.Require().Equal(ts.newText, ts.db.GetTweet(strconv.Itoa(tweet1.ID)).Text)
}

func (ts *MapStoreTestSuit) TestDeleteTweet() {
	createdUser1 := ts.db.AddUser(&ts.user1)
	tweet1 := ts.db.AddTweet(&ts.tweet1, createdUser1)
	tweets := ts.db.GetUserTweets(strconv.Itoa(createdUser1.ID))
	countOfTweets := 1
	ts.Require().Equal(countOfTweets, len(tweets))
	isDeleted := ts.db.DeleteTweet(strconv.Itoa(tweet1.ID), createdUser1)
	ts.Require().True(isDeleted)
	tweets = ts.db.GetUserTweets(strconv.Itoa(createdUser1.ID))
	countOfTweets = 0
	ts.Require().Equal(countOfTweets, len(tweets))

}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MapStoreTestSuit))
}
