package main

import "github.com/go-pg/pg"

const (
	host     = "db"
	port     = 5432
	user     = "admin"
	password = "password"
	dbname   = "twitter"
)

var pgOptions = pg.Options{User: user, Password: password, Database: dbname}
