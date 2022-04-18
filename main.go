package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Configuration struct {
	Port     int
	MaxJokes int
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Joke struct {
	Id       string `json:"id"`
	Icon_url string `json:"icon_url"`
	Value    string `json:"value"`
}

type DB struct {
	jokes []Joke
}

func (config *Configuration) load() error {
	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		return errors.New("PORT is not specified or is not a positive integer")
	}

	maxJokes, err := strconv.Atoi(os.Getenv("MAX_JOKES"))

	if err != nil {
		return errors.New("MAX_JOKES is not specified or is not a positive integer")
	}

	config.Port = port
	config.MaxJokes = maxJokes
	return nil
}

func (config *Configuration) portToString() string {
	return fmt.Sprintf(":%d", config.Port)
}

func (db *DB) clean() {
	db.jokes = []Joke{}
}

func (db *DB) addJoke() {
	exists := false
	j, _ := getJoke()
	for _, joke := range db.jokes {
		if joke.Id == j.Id {
			exists = true
		}
	}
	if !exists && len(db.jokes) < config.MaxJokes {
		db.jokes = append(db.jokes, j)
	}
}

func (db *DB) addJokeWG(wg *sync.WaitGroup) {
	defer wg.Done()
	db.addJoke()
}

func (db *DB) addJokeChanel(c chan<- int) {
	db.addJoke()
	c <- 1
}

func getJoke() (Joke, error) {
	var j Joke
	resp, _ := http.Get("https://api.chucknorris.io/jokes/random")
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &j)
	return j, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	msg := MessageResponse{"Jokes API"}
	js, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	msg := MessageResponse{"pong"}
	js, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func jokesSyncHandler(w http.ResponseWriter, r *http.Request) {
	db.clean()

	for len(db.jokes) < config.MaxJokes {
		db.addJoke()
	}

	js, _ := json.Marshal(db.jokes)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func jokesWGHandler(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	db.clean()

	for len(db.jokes) < config.MaxJokes {
		free := config.MaxJokes - len(db.jokes)
		for i := 0; i < free; i++ {
			wg.Add(1)
			go db.addJokeWG(&wg)
		}
		wg.Wait()
	}

	js, _ := json.Marshal(db.jokes)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func jokesChanelHandler(w http.ResponseWriter, r *http.Request) {
	c := make(chan int, config.MaxJokes)
	db.clean()

	for len(db.jokes) < config.MaxJokes {
		free := config.MaxJokes - len(db.jokes)
		for i := 0; i < free; i++ {
			go db.addJokeChanel(c)
		}

		for i := 0; i < free; i++ {
			<-c
		}
	}

	js, _ := json.Marshal(db.jokes)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

var config Configuration
var db DB

func main() {

	err := config.load()

	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/jokes/sync", jokesSyncHandler)
	http.HandleFunc("/jokes/wg", jokesWGHandler)
	http.HandleFunc("/jokes/chanel", jokesChanelHandler)

	err = http.ListenAndServe(config.portToString(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
