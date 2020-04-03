package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Alert struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID       int    `json:"id",pg:",unique"`
	Login    string `json:"login",pg:",unique"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type Tweet struct {
	ID       int       `json:"id"`
	Time     time.Time `json:"time"`
	Author   string    `json:"author"`
	AuthorID int       `json:"authorID",pg:"fk:user_id"`
	Text     string    `json:"text"`
}

type Tok struct {
	TokenValue string `json:"token"`
}

type jwtUserClaim struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	jwt.StandardClaims
}
