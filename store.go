package main

import (
	"sort"
	"strconv"
	"sync"
	"time"
)

type Store interface {
	getTweet(id string) *Tweet
	getTweets() *[]Tweet
	addTweet(tweet *Tweet, user *User) *Tweet
	updateTweet(tweetID string, user User, tweet Tweet) bool
	deleteTweet(tweetID string, user User) bool
	addUser(user User) *User
	getUserByID(userID string) *User
	checkLoginPassword(login string, password string) bool
	getUserTweets(authorID string) []*Tweet
}

type MapStore struct {
	users   []User
	tweets  map[string][]*Tweet
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

func (ms *MapStore) getTweets() []*Tweet {
	var allTweets []*Tweet
	for _, us := range ms.users {
		allTweets = append(allTweets, ms.tweets[us.ID]...)
	}
	sort.Sort(byID(allTweets))
	return allTweets
}

func (ms *MapStore) deleteTweet(tweetID string, user *User) bool {
	for index, twt := range ms.tweets[user.ID] {
		if twt.ID == tweetID {
			twts := append(ms.tweets[user.ID][:index], ms.tweets[user.ID][index+1:]...)
			ms.tweets[user.ID] = twts
			return true
		}
	}
	return false
}

func (ms *MapStore) addUser(user *User) *User {
	user.ID = strconv.Itoa(ms.getUserID())
	ms.users = append(ms.users, *user)
	return user
}

func (ms *MapStore) addTweet(tweet *Tweet, user *User) *Tweet {
	tweet.ID = strconv.Itoa(ms.getTweetID())
	tweet.AuthorID = user.ID
	ms.tweets[user.ID] = append(ms.tweets[user.ID], tweet)
	return tweet
}

func (ms *MapStore) getTweet(id string) *Tweet {
	for us := range ms.tweets {
		for _, twt := range ms.tweets[us] {
			if twt.ID == id {
				return twt
			}
		}
	}
	return nil
}

func (ms *MapStore) updateTweet(tweetID string, user *User, text string) bool {
	for index, twt := range ms.tweets[user.ID] {
		if twt.ID == tweetID {
			changeTweet := ms.tweets[user.ID][index]
			changeTweet.Text = text
			changeTweet.Time = time.Now()

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

func (ms *MapStore) getUserTweets(authorID string) []*Tweet {
	userTweets := ms.tweets[authorID]
	sort.Sort(byID(userTweets))
	return userTweets
}

func (ms *MapStore) getTweetID() int {
	ms.mutex.Lock()
	ms.tweetID += 1
	ms.mutex.Unlock()
	return ms.tweetID
}

func (ms *MapStore) getUserID() int {
	ms.userID += 1
	return ms.userID
}
