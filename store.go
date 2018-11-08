package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	redis "github.com/go-redis/redis"
)

func initRedisConnection() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getCachedMeetupName(memberID int, tokenString string, deck *DB_Deck) string {
	// memberID string, eventID string, tokenString string, groupname string)
	card, err := getCard(memberID)
	if err != nil {
		cards := getMeetupEventMembersArray(tokenString, deck.GroupName, deck.EventID)
		for i := 0; i < len(cards); i++ {
			fmt.Println(cards[i].MeetupID)
			err := cacheCard(cards[i].MeetupID, cards[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	card, err = getCard(memberID)
	if err != nil {
		log.Fatal(err)
	}
	return card.Name
}

func getCachedMeetupImage(memberID int, tokenString string, deck *DB_Deck) string {
	// memberID string, eventID string, tokenString string, groupname string)
	card, err := getCard(memberID)
	if err != nil {
		cards := getMeetupEventMembersArray(tokenString, deck.GroupName, deck.EventID)
		for i := 0; i < len(cards); i++ {
			err := cacheCard(cards[i].MeetupID, cards[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	card, err = getCard(memberID)
	if err != nil {
		log.Fatal(err)
	}
	return card.ImageURL
}

func getCard(memberID int) (Card, error) {
	res, err := redisClient.Get(strconv.Itoa(memberID)).Result()
	var cachedCard Card
	if err != nil {
		return cachedCard, err
	}
	err = json.Unmarshal([]byte(res), &cachedCard)
	return cachedCard, err
}

func cacheCard(memberID int, card Card) error {
	expiration, err := time.ParseDuration("27h")
	if err != nil {
		fmt.Println(err)
	}
	redisObj, err := json.Marshal(card)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(card)
	// fmt.Println(memberID)
	err = redisClient.Set(strconv.Itoa(memberID), redisObj, expiration).Err()
	return err
}
