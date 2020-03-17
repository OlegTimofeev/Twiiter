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
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type Tweet struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	Author   string    `json:"author"`
	AuthorID string    `json:"authorID"`
	Text     string    `json:"text"`
}

type jwtUserClaim struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	jwt.StandardClaims
}
