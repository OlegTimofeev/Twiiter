package main

import (
	"github.com/go-pg/pg"
	"sync"
)

type Store interface {
	InitDB() error
	GetTweet(id string) (*Tweet, error)
	GetTweets() ([]*Tweet, error)
	AddTweet(tweet *Tweet, user *User) (*Tweet, error)
	UpdateTweet(tweetID string, user *User, text string) (bool, error)
	DeleteTweet(tweetID string, user *User) (bool, error)
	AddUser(user User) (*User, error)
	GetUserByID(userID string) (*User, error)
	CheckLoginPassword(login string, password string) (*User, bool, error)
	GetUserTweets(authorID string) ([]*Tweet, error)
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
