package main

import (
	"sort"
	"strconv"
	"time"
)

func (ms *MapStore) GetUserByID(userID string) *User {
	for _, item := range ms.users {
		if strconv.Itoa(item.ID) == userID {
			return &item
		}
	}
	return nil
}

func (ms *MapStore) GetTweets() []*Tweet {
	var allTweets []*Tweet
	for _, us := range ms.users {
		allTweets = append(allTweets, ms.tweets[us.ID]...)
	}
	sort.Sort(byID(allTweets))
	return allTweets
}

func (ms *MapStore) DeleteTweet(tweetID string, user *User) bool {
	for index, twt := range ms.tweets[user.ID] {
		if strconv.Itoa(twt.ID) == tweetID {
			twts := append(ms.tweets[user.ID][:index], ms.tweets[user.ID][index+1:]...)
			ms.tweets[user.ID] = twts
			return true
		}
	}
	return false
}

func (ms *MapStore) AddUser(user *User) *User {
	user.ID = ms.GetUserID()
	ms.users = append(ms.users, *user)
	return user
}

func (ms *MapStore) AddTweet(tweet *Tweet, user *User) *Tweet {
	tweet.ID = ms.GetTweetID()
	tweet.AuthorID = user.ID
	ms.tweets[user.ID] = append(ms.tweets[user.ID], tweet)
	return tweet
}

func (ms *MapStore) GetTweet(id string) *Tweet {
	for us := range ms.tweets {
		for _, twt := range ms.tweets[us] {
			if strconv.Itoa(twt.ID) == id {
				return twt
			}
		}
	}
	return nil
}

func (ms *MapStore) UpdateTweet(tweetID string, user *User, text string) bool {
	for index, twt := range ms.tweets[user.ID] {
		if strconv.Itoa(twt.ID) == tweetID {
			changeTweet := ms.tweets[user.ID][index]
			changeTweet.Text = text
			changeTweet.Time = time.Now()

			return true
		}
	}
	return false
}

func (ms *MapStore) CheckLoginPassword(login string, password string) (*User, bool) {
	for _, bdUser := range ms.users {
		if login == bdUser.Login && password == bdUser.Password {
			return &bdUser, true
		}
	}
	return nil, false
}

func (ms *MapStore) GetUserTweets(authorID string) []*Tweet {
	id, _ := strconv.Atoi(authorID)
	userTweets := ms.tweets[id]
	sort.Sort(byID(userTweets))
	return userTweets
}

func (ms *MapStore) GetTweetID() int {
	ms.mutex.Lock()
	ms.tweetID += 1
	ms.mutex.Unlock()
	return ms.tweetID
}

func (ms *MapStore) GetUserID() int {
	ms.userID += 1
	return ms.userID
}
