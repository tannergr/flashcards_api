package main

import (
	"database/sql"
	_ "html/template"
	_ "io/ioutil"
	"log"
	"net/http"
	_ "net/url"
	"os"

	redis "github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB
var redisClient *redis.Client

func main() {

	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	redisClient, err = initRedisConnection()
	if err != nil {
		log.Fatal("Could not connect to redis")
	}

	r := mux.NewRouter()

	r.Handle("/", IndexHandler).Methods("GET")
	r.Handle("/test", TestLimit).Methods("GET")
	r.Handle("/api/meetup/auth", CallBackHandler).Methods("GET")
	r.Handle("/api/events", GetEvents).Methods("GET")
	r.Handle("/api/events/{groupname}/{eid}", getEventMembers).Methods("GET")
	r.Handle("/api/member", GetMember).Methods("GET")
	r.Handle("/api/deck", AddDeck).Methods("POST")
	r.Handle("/api/deck", GetDecks).Methods("GET")
	r.Handle("/api/deck/{deckID}", GetCards).Methods("GET")
	r.Handle("/api/deck/{deckID}/score/{score}", PostScore).Methods("POST")
	r.Handle("/api/deck/{deckID}/select", SetSelectedDeck).Methods("Post")
	r.Handle("/api/member/deck", GetLastDeck).Methods("GET")

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authentication", "Content-Type"},
		// Debug:            true,
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handler)))
}
