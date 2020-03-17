package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func (ts *MapStoreTestSuit) SetupTest() {
	ts.db = &MapStore{tweets: make(map[string][]Tweet), userID: 0, tweetID: 0}
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

	createdUser := ts.db.addUser(ts.user1)
	us, isExist := ts.db.checkLoginPassword(createdUser.Login, createdUser.Password)
	assert.True(ts.T(), isExist)
	assert.Equal(ts.T(), us.Login, createdUser.Login)
	assert.Equal(ts.T(), us.Password, createdUser.Password)
	assert.Equal(ts.T(), us.Name, createdUser.Name)
	assert.Equal(ts.T(), us.Surname, createdUser.Surname)
}

func (ts *MapStoreTestSuit) TestAddTweet() {
	createdUser := ts.db.addUser(ts.user1)
	tweet := ts.db.addTweet(ts.tweet1, *createdUser)
	assert.Equal(ts.T(), *tweet, *ts.db.getTweet(tweet.ID))
}

func (ts *MapStoreTestSuit) TestGetUserTweets() {
	createdUser1 := ts.db.addUser(ts.user1)
	assert.Nil(ts.T(), *ts.db.getUserTweets(createdUser1.ID))
	createdUser2 := ts.db.addUser(ts.user2)
	tweet1 := ts.db.addTweet(ts.tweet1, *createdUser1)
	tweets := ts.db.getUserTweets(createdUser1.ID)
	assert.Contains(ts.T(), *tweets, *tweet1)
	tweet2 := ts.db.addTweet(ts.tweet1, *createdUser1)
	tweet3 := ts.db.addTweet(ts.tweet3, *createdUser2)
	tweets = ts.db.getUserTweets(createdUser1.ID)
	countOfTweets := 2
	assert.NotContains(ts.T(), *tweets, *tweet3)
	assert.Contains(ts.T(), *tweets, *tweet1)
	assert.Contains(ts.T(), *tweets, *tweet2)
	assert.Equal(ts.T(), countOfTweets, len(*tweets))
}

func (ts *MapStoreTestSuit) TestUpdateTweet() {
	createdUser1 := ts.db.addUser(ts.user1)
	tweet1 := ts.db.addTweet(ts.tweet1, *createdUser1)
	isChanged := ts.db.updateTweet(tweet1.ID, *createdUser1, ts.newText)
	assert.True(ts.T(), isChanged)
	assert.Equal(ts.T(), ts.newText, ts.db.getTweet(tweet1.ID).Text)
}

func (ts *MapStoreTestSuit) TestDeleteTweet() {
	createdUser1 := ts.db.addUser(ts.user1)
	tweet1 := ts.db.addTweet(ts.tweet1, *createdUser1)
	tweets := ts.db.getUserTweets(createdUser1.ID)
	countOfTweets := 1
	assert.Equal(ts.T(), countOfTweets, len(*tweets))
	isDeleted := ts.db.deleteTweet(tweet1.ID, *createdUser1)
	assert.True(ts.T(), isDeleted)
	tweets = ts.db.getUserTweets(createdUser1.ID)
	countOfTweets = 0
	assert.Equal(ts.T(), countOfTweets, len(*tweets))

}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MapStoreTestSuit))
}
