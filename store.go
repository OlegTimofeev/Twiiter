package main

import (
	"sort"
	"time"
)

type Store interface {
	getTweet(id string) *Tweet
	getTweets() []Tweet
	addTweet(tweet Tweet, user User)
	updateTweet(tweetId string, userId string, tweet Tweet) bool
	deleteTweet(tweetId string, userId string) bool
	addUser(user User)
	getUserById(userId string) User
	checkLoginPassword(login string, password string) bool
}

type mapa struct {
	users  []User
	tweets map[string][]Tweet
}

func (mapa) getUserById(userId string) *User {
	for _, item := range users {
		if item.ID == userId {
			return &item
		}
	}
	return nil
}

func (mapa) getTweets() []Tweet {
	var allTweets []Tweet
	for _, us := range users {
		allTweets = append(allTweets, tweets[us.ID]...)
	}
	sort.Sort(byId(allTweets))
	return allTweets
}

func (mapa) deleteTweet(tweetId string, userId string) bool {
	for index, twt := range tweets[userId] {
		if twt.ID == tweetId {
			twts := append(tweets[userId][:index], tweets[userId][index+1:]...)
			tweets[userId] = twts
			return true
		}
	}
	return false
}

func (mapa) addUser(user User) {
	users = append(users, user)
}

func (mapa) addTweet(tweet Tweet, user User) {
	tweets[user.ID] = append(tweets[user.ID], tweet)
}

func (mapa) getTweet(id string) *Tweet {
	for us := range tweets {
		for _, twt := range tweets[us] {
			if twt.ID == id {
				return &twt
			}
		}
	}
	return nil
}

func (mapa) updateTweet(tweetId string, userId string, text string) bool {
	for index, twt := range tweets[userId] {
		if twt.ID == tweetId {
			changeTweet := &tweets[userId][index]
			changeTweet.Text = text
			changeTweet.Time = time.Now().Format("2006-01-02 15:04")

			return true
		}
	}
	return false
}

func (mapa) checkLoginPassword(login string, password string) (*User, bool) {
	for _, bdUser := range users {
		if login == bdUser.Login && password == bdUser.Password {
			return &bdUser, true
		}
	}
	return nil, false
}
