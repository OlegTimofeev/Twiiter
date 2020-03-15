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
	getUserTweets(authorId string) []Tweet
}

type MapStore struct {
	users  []User
	tweets map[string][]Tweet
}

func (mapStore *MapStore) getUserById(userId string) *User {
	for _, item := range mapStore.users {
		if item.ID == userId {
			return &item
		}
	}
	return nil
}

func (mapStore *MapStore) getTweets() []Tweet {
	var allTweets []Tweet
	for _, us := range mapStore.users {
		allTweets = append(allTweets, mapStore.tweets[us.ID]...)
	}
	sort.Sort(byId(allTweets))
	return allTweets
}

func (mapStore *MapStore) deleteTweet(tweetId string, userId string) bool {
	for index, twt := range mapStore.tweets[userId] {
		if twt.ID == tweetId {
			twts := append(mapStore.tweets[userId][:index], mapStore.tweets[userId][index+1:]...)
			mapStore.tweets[userId] = twts
			return true
		}
	}
	return false
}

func (mapStore *MapStore) addUser(user User) {
	mapStore.users = append(mapStore.users, user)
}

func (mapStore *MapStore) addTweet(tweet Tweet, user User) {
	mapStore.tweets[user.ID] = append(mapStore.tweets[user.ID], tweet)
}

func (mapStore *MapStore) getTweet(id string) *Tweet {
	for us := range mapStore.tweets {
		for _, twt := range mapStore.tweets[us] {
			if twt.ID == id {
				return &twt
			}
		}
	}
	return nil
}

func (mapStore *MapStore) updateTweet(tweetId string, userId string, text string) bool {
	for index, twt := range mapStore.tweets[userId] {
		if twt.ID == tweetId {
			changeTweet := &mapStore.tweets[userId][index]
			changeTweet.Text = text
			changeTweet.Time = time.Now().Format("2006-01-02 15:04")

			return true
		}
	}
	return false
}

func (mapStore *MapStore) checkLoginPassword(login string, password string) (*User, bool) {
	for _, bdUser := range mapStore.users {
		if login == bdUser.Login && password == bdUser.Password {
			return &bdUser, true
		}
	}
	return nil, false
}

func (mapStore *MapStore) getUserTweets(authorId string) []Tweet {
	userTweets := mapStore.tweets[authorId]
	sort.Sort(byId(userTweets))
	return userTweets
}
