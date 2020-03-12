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
	w.Header().Set("Content-Type", "application/json")
	var tweet Tweet
	_ = json.NewDecoder(r.Body).Decode(&tweet)
	tweet.ID = strconv.Itoa(rand.Intn(10000))
	tweet.Time = time.Now().Format("2006-01-02 15:04")
	tweets = append(tweets, tweet)
	json.NewEncoder(w).Encode(tweet)
}

func deleteTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range tweets {
		if item.ID == params["id"] {
			tweets = append(tweets[:index], tweets[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tweets)
}

func updateTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range tweets {
		if item.ID == params["id"] {
			tweets = append(tweets[:index], tweets[index+1:]...)
			var tweet Tweet
			_ = json.NewDecoder(r.Body).Decode(&tweet)
			tweet.ID = params["id"]
			tweet.Time = time.Now().Format("2006-01-02 15:04")
			tweets = append(tweets, tweet)
			json.NewEncoder(w).Encode(tweet)
			return
		}
	}
	json.NewEncoder(w).Encode(tweets)
}

func main() {
	r := mux.NewRouter()
	tweets = append(tweets, Tweet{ID: "1", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u", Author: "Daria"})
	tweets = append(tweets, Tweet{ID: "2", Time: time.Now().Format("2006-01-02 15:04"), Text: "I love u 2 Daria", Author: "Oleg"})
	r.HandleFunc("/tweets", getTweets).Methods("GET")
	r.HandleFunc("/tweets/{id}", getTweet).Methods("GET")
	r.HandleFunc("/tweets", createTweet).Methods("POST")
	r.HandleFunc("/tweets/{id}", updateTweet).Methods("PUT")
	r.HandleFunc("/tweets/{id}", deleteTweet).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
