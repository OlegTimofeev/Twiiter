package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Tweet struct {
	ID     string `json:"id"`
	Time   string `json:"time"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

var tweets []Tweet
var users []User
var currentuser = new(User)
var error = Error{Name: "Not Auth", Description: "Auth to create,update or delete twit"}

func getTweets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}

func getTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range tweets {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Tweet{})
}

func createTweet(w http.ResponseWriter, r *http.Request) {
	if currentuser.ID != "" {
		w.Header().Set("Content-Type", "application/json")
		var tweet Tweet
		_ = json.NewDecoder(r.Body).Decode(&tweet)
		tweet.ID = strconv.Itoa(rand.Intn(10000))
		tweet.Time = time.Now().Format("2006-01-02 15:04")
		tweet.Author = currentuser.Name + " " + currentuser.Surname
		tweets = append(tweets, tweet)
		json.NewEncoder(w).Encode(tweet)
	} else {
		json.NewEncoder(w).Encode(error)
	}

}

func deleteTweet(w http.ResponseWriter, r *http.Request) {
	if currentuser.ID != "" {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range tweets {
			if item.ID == params["id"] {
				tweets = append(tweets[:index], tweets[index+1:]...)
				break
			}
		}
		json.NewEncoder(w).Encode(tweets)
	} else {
		json.NewEncoder(w).Encode(error)
	}
}

func updateTweet(w http.ResponseWriter, r *http.Request) {
	if currentuser.ID != "" {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range tweets {
			if item.ID == params["id"] {
				tweets = append(tweets[:index], tweets[index+1:]...)
				var tweet Tweet
				_ = json.NewDecoder(r.Body).Decode(&tweet)
				tweet.ID = params["id"]
				tweet.Author = currentuser.Name + " " + currentuser.Surname
				tweet.Time = time.Now().Format("2006-01-02 15:04")
				tweets = append(tweets, tweet)
				json.NewEncoder(w).Encode(tweet)
				return
			}
		}
		json.NewEncoder(w).Encode(tweets)
	} else {
		json.NewEncoder(w).Encode(error)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentuser.ID = strconv.Itoa(rand.Intn(1000))
	_ = json.NewDecoder(r.Body).Decode(currentuser)
	users = append(users, *currentuser)
	json.NewEncoder(w).Encode(currentuser)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if currentuser.ID == "" {
		var loginUser User
		_ = json.NewDecoder(r.Body).Decode(&loginUser)
		for _, bdUser := range users {
			if loginUser.Login == bdUser.Login && loginUser.Password == bdUser.Password {
				currentuser = &bdUser
				json.NewEncoder(w).Encode(currentuser)
				return
			}
		}
	}
	json.NewEncoder(w).Encode(tweets)
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentuser = new(User)
	json.NewEncoder(w).Encode(tweets)
}

func main() {
	r := mux.NewRouter()
	tweets = append(tweets, Tweet{ID: "1", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u", Author: "Daria"})
	tweets = append(tweets, Tweet{ID: "2", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u 2 Daria", Author: "Oleg"})
	users = append(users, User{ID: "3", Login: "www", Password: "123", Name: "Ol", Surname: "eg"})
	r.HandleFunc("/tweets", getTweets).Methods("GET")
	r.HandleFunc("/tweets/{id}", getTweet).Methods("GET")
	r.HandleFunc("/tweets", createTweet).Methods("POST")
	r.HandleFunc("/tweets/{id}", updateTweet).Methods("PUT")
	r.HandleFunc("/tweets/{id}", deleteTweet).Methods("DELETE")
	r.HandleFunc("/tweets/signUp", signUp).Methods("POST")
	r.HandleFunc("/tweets/signIn", signIn).Methods("POST")
	r.HandleFunc("/tweets/logout", logout).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
