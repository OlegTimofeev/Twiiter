package main

import (
	"github.com/go-pg/pg"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Store interface {
	getTweet(id string) *Tweet
	getTweets() []*Tweet
	addTweet(tweet *Tweet, user *User) *Tweet
	updateTweet(tweetID string, user User, text string) bool
	deleteTweet(tweetID string, user User) bool
	addUser(user User) *User
	getUserByID(userID string) *User
	checkLoginPassword(login string, password string) bool
	getUserTweets(authorID string) []*Tweet
}

type MapStore struct {
	users   []User
	tweets  map[int][]*Tweet
	mutex   sync.Mutex
	userID  int
	tweetID int
}

type PostgresDB struct {
	pgdb *pg.DB
}

func (ms *MapStore) getUserByID(userID string) *User {
	for _, item := range ms.users {
		if strconv.Itoa(item.ID) == userID {
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
		if strconv.Itoa(twt.ID) == tweetID {
			twts := append(ms.tweets[user.ID][:index], ms.tweets[user.ID][index+1:]...)
			ms.tweets[user.ID] = twts
			return true
		}
	}
	return false
}

func (ms *MapStore) addUser(user *User) *User {
	user.ID = ms.getUserID()
	ms.users = append(ms.users, *user)
	return user
}

func (ms *MapStore) addTweet(tweet *Tweet, user *User) *Tweet {
	tweet.ID = ms.getTweetID()
	tweet.AuthorID = user.ID
	ms.tweets[user.ID] = append(ms.tweets[user.ID], tweet)
	return tweet
}

func (ms *MapStore) getTweet(id string) *Tweet {
	for us := range ms.tweets {
		for _, twt := range ms.tweets[us] {
			if strconv.Itoa(twt.ID) == id {
				return twt
			}
		}
	}
	return nil
}

func (ms *MapStore) updateTweet(tweetID string, user *User, text string) bool {
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

func (ms *MapStore) checkLoginPassword(login string, password string) (*User, bool) {
	for _, bdUser := range ms.users {
		if login == bdUser.Login && password == bdUser.Password {
			return &bdUser, true
		}
	}
	return nil, false
}

func (ms *MapStore) getUserTweets(authorID string) []*Tweet {
	id, _ := strconv.Atoi(authorID)
	userTweets := ms.tweets[id]
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

func (dbpg PostgresDB) addTweet(tweet *Tweet, user *User) *Tweet {
	tweet.AuthorID = user.ID
	dbpg.pgdb.Insert(tweet)
	return tweet
}

func (dbpg PostgresDB) checkLoginPassword(login string, password string) (*User, bool) {
	user := new(User)
	if err := dbpg.pgdb.Model(user).Where("login = ?", login).Where("password = ?", password).Select(); err != nil {
		return nil, false
	}
	return user, true
}

func (dbpg PostgresDB) addUser(user *User) *User {
	if err := dbpg.pgdb.Insert(user); err != nil {
		return nil
	}
	return user
}

func (dbpg PostgresDB) getUserByID(userID string) *User {
	id, _ := strconv.Atoi(userID)
	user := User{ID: id}
	if err := dbpg.pgdb.Select(&user); err != nil {
		return nil
	}
	return &user
}

func (dbpg PostgresDB) deleteTweet(tweetID string, user User) bool {
	twt := dbpg.getTweetCheckAuthor(tweetID, user)
	if twt == nil {
		return false
	}
	if err := dbpg.pgdb.Delete(twt); err != nil {
		return false
	}
	return true
}

func (dbpg PostgresDB) updateTweet(tweetID string, user User, text string) bool {
	twt := dbpg.getTweetCheckAuthor(tweetID, user)
	if twt == nil {
		return false
	}
	twt.Text = text
	if err := dbpg.pgdb.Update(twt); err != nil {
		return false
	}
	return true
}

func (dbpg PostgresDB) getTweetCheckAuthor(tweetID string, user User) *Tweet {
	twtID, _ := strconv.Atoi(tweetID)
	twt := Tweet{ID: twtID}
	if err := dbpg.pgdb.Select(&twt); err != nil {
		return nil
	}
	if twt.AuthorID != user.ID {
		return nil
	}
	return &twt
}

func (dbpg PostgresDB) getTweet(id string) *Tweet {
	twtID, _ := strconv.Atoi(id)
	twt := Tweet{ID: twtID}
	if err := dbpg.pgdb.Select(&twt); err != nil {
		return nil
	}
	return &twt
}

func (dbpg PostgresDB) getTweets() []*Tweet {
	var tweets []*Tweet
	if err := dbpg.pgdb.Model(&tweets).Select(); err != nil {
		return nil
	}
	return tweets
}

func (dbpg PostgresDB) getUserTweets(authorID string) []*Tweet {
	id, _ := strconv.Atoi(authorID)
	var tweets []*Tweet
	if err := dbpg.pgdb.Model(&tweets).Where("author_id = ?", id).Select(); err != nil {
		return nil
	}
	return tweets
}
