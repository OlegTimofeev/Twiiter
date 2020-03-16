package main

import (
	"sort"
	"strconv"
	"sync"
	"time"
)

type Store interface {
	getTweet(id string) *Tweet
	getTweets() []Tweet
	addTweet(tweet Tweet, user User)
	updateTweet(tweetID string, userID string, tweet Tweet) bool
	deleteTweet(tweetID string, userID string) bool
	addUser(user User) User
	getUserByID(userID string) User
	checkLoginPassword(login string, password string) bool
	getUserTweets(authorID string) []Tweet
}

type MapStore struct {
	users   []User
	tweets  map[string][]Tweet
	mutex   sync.Mutex
	userID  int
	tweetID int
}

func (ms *MapStore) getUserByID(userID string) *User {
	for _, item := range ms.users {
		if item.ID == userID {
			return &item
		}
	}
	return nil
}

func (ms *MapStore) getTweets() []Tweet {
	var allTweets []Tweet
	for _, us := range ms.users {
		allTweets = append(allTweets, ms.tweets[us.ID]...)
	}
	sort.Sort(byID(allTweets))
	return allTweets
}

func (ms *MapStore) deleteTweet(tweetID string, userID string) bool {
	for index, twt := range ms.tweets[userID] {
		if twt.ID == tweetID {
			twts := append(ms.tweets[userID][:index], ms.tweets[userID][index+1:]...)
			ms.tweets[userID] = twts
			return true
		}
	}
	return false
}

func (ms *MapStore) addUser(user User) User {
	ms.mutex.Lock()
	user.ID = strconv.Itoa(ms.getUserID())
	ms.mutex.Unlock()
	ms.users = append(ms.users, user)
	return user
}

func (ms *MapStore) addTweet(tweet Tweet, user User) Tweet {
	ms.mutex.Lock()
	tweet.ID = strconv.Itoa(ms.getTweetID())
	ms.mutex.Unlock()
	tweet.AuthorID = user.ID
	ms.tweets[user.ID] = append(ms.tweets[user.ID], tweet)
	return tweet
}

func (ms *MapStore) getTweet(id string) *Tweet {
	for us := range ms.tweets {
		for _, twt := range ms.tweets[us] {
			if twt.ID == id {
				return &twt
			}
		}
	}
	return nil
}

func (ms *MapStore) updateTweet(tweetID string, userID string, text string) bool {
	for index, twt := range ms.tweets[userID] {
		if twt.ID == tweetID {
			changeTweet := &ms.tweets[userID][index]
			changeTweet.Text = text
			changeTweet.Time = time.Now().Format("2006-01-02 15:04")

			return true
		}
	}
	return false
}

func (ms *MapStore) checkLoginPassword(login string, password string) (*User, bool) {
	for _, bdUser := range ms.users {
		if login == bdUser.Login && password == bdUser.Password {
			return &bdUser, true
		}
	}
	return nil, false
}

func (ms *MapStore) getUserTweets(authorID string) []Tweet {
	userTweets := ms.tweets[authorID]
	sort.Sort(byID(userTweets))
	return userTweets
}

func (ms *MapStore) getTweetID() int {
	ms.tweetID += 1
	return ms.tweetID
}

func (ms *MapStore) getUserID() int {
	ms.userID += 1
	return ms.userID
}
