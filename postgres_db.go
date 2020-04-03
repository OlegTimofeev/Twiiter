package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"strconv"
)

func connect() *pg.DB {
	return pg.Connect(&pgOptions)
}

func (db *PostgresDB) InitDB() error {
	db.pgdb = connect()
	err := db.pgdb.RunInTransaction(func(tx *pg.Tx) error {
		panicIf(db.pgdb.DropTable((*User)(nil), &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		}))
		panicIf(db.pgdb.CreateTable((*User)(nil), &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}))
		panicIf(db.pgdb.DropTable((*Tweet)(nil), &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		}))
		panicIf(db.pgdb.CreateTable((*Tweet)(nil), &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}))
		if cons := db.pgdb.PoolStats().Hits; cons < 1 {
			return *new(error)
		}
		return nil
	})
	return err
}

func (db *PostgresDB) AddTweet(tweet *Tweet, user *User) (*Tweet, error) {
	tweet.AuthorID = user.ID
	if err := db.pgdb.Insert(tweet); err != nil {
		return nil, err
	}
	return tweet, nil
}

func (db *PostgresDB) CheckLoginPassword(login string, password string) (*User, bool, error) {
	user := new(User)
	if err := db.pgdb.Model(user).Where("login = ?", login).Where("password = ?", password).Select(); err != nil {
		return nil, false, err
	}
	return user, true, nil
}

func (db *PostgresDB) AddUser(user *User) (*User, error) {
	if err := db.pgdb.Insert(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (db *PostgresDB) GetUserByID(userID string) (*User, error) {
	id, _ := strconv.Atoi(userID)
	user := User{ID: id}
	if err := db.pgdb.Select(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *PostgresDB) DeleteTweet(tweetID string, user *User) (bool, error) {
	err := db.pgdb.RunInTransaction(func(tx *pg.Tx) error {
		twt, err := db.GetTweetCheckAuthor(tweetID, user)
		if twt == nil {
			return err
		}
		if err = db.pgdb.Delete(twt); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *PostgresDB) UpdateTweet(tweetID string, user *User, text string) (bool, error) {
	err := db.pgdb.RunInTransaction(func(tx *pg.Tx) error {
		twt, err := db.GetTweetCheckAuthor(tweetID, user)
		if twt == nil {
			return err
		}
		twt.Text = text
		if err := db.pgdb.Update(twt); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *PostgresDB) GetTweetCheckAuthor(tweetID string, user *User) (*Tweet, error) {
	twtID, _ := strconv.Atoi(tweetID)
	twt := Tweet{ID: twtID}
	if err := db.pgdb.Select(&twt); err != nil {
		return nil, err
	}
	if twt.AuthorID != user.ID {
		return nil, nil
	}
	return &twt, nil
}

func (db *PostgresDB) GetTweet(id string) (*Tweet, error) {
	twtID, _ := strconv.Atoi(id)
	twt := Tweet{ID: twtID}
	if err := db.pgdb.Select(&twt); err != nil {
		return nil, err
	}
	return &twt, nil
}

func (db *PostgresDB) GetTweets() ([]*Tweet, error) {
	var tweets []*Tweet
	if err := db.pgdb.Model(&tweets).Select(); err != nil {
		return nil, err
	}
	return tweets, nil
}

func (db *PostgresDB) GetUserTweets(authorID string) ([]*Tweet, error) {
	id, _ := strconv.Atoi(authorID)
	var tweets []*Tweet
	if err := db.pgdb.Model(&tweets).Where("author_id = ?", id).Select(); err != nil {
		return nil, err
	}
	return tweets, nil
}
