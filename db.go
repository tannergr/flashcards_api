package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	connString, err := ioutil.ReadFile("./config/pgConnectionString")

	if err != nil {
		log.Fatal("Error reading private key")
		return nil, err
	}

	return sql.Open("postgres", string(connString))
}

func addDeckToDB(memberID int64, deck Deck) error {
	fmt.Printf("%v : %v", memberID, deck.Name)

	query := `
		INSERT INTO decks (deck_name, user_id, event_id)
		VALUES($1, $2, $3)
		RETURNING deck_id
	`
	deckID := 0
	err := db.QueryRow(query, deck.Name, memberID, deck.EventID).Scan(&deckID)
	if err != nil {
		return err
	}

	// https://stackoverflow.com/questions/12486436/golang-how-do-i-batch-sql-statements-with-package-database-sql/25192138#25192138
	// Batch insert
	valueStrings := make([]string, 0, len(deck.Cards))
	valueArgs := make([]interface{}, 0, len(deck.Cards)*3)
	i := 1
	for _, card := range deck.Cards {
		valueString := fmt.Sprintf("($%v, $%v, $%v)", i, i+1, i+2)
		valueStrings = append(valueStrings, valueString)
		valueArgs = append(valueArgs, deckID)
		valueArgs = append(valueArgs, card.Name)
		valueArgs = append(valueArgs, card.ImageURL)
		i += 3
	}

	query = fmt.Sprintf(`
		WITH rows AS (
			INSERT INTO members (deck_id, meetup_id, image_url)
			VALUES %v
			RETURNING 1
		)
		SELECT count(*) FROM rows;
	`, strings.Join(valueStrings, ","))

	count := 0
	err = db.QueryRow(query, valueArgs...).Scan(&count)
	if err != nil {
		return err
	}
	fmt.Println(err)

	return err
}

func getDecksFromDB(memberID int64) ([]*DB_Deck, error) {
	query := `
		SELECT * from DECKS
		WHERE user_id=$1
	`

	rows, err := db.Query(query, memberID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	decks := make([]*DB_Deck, 0)
	for rows.Next() {
		deck := new(DB_Deck)
		err := rows.Scan(&deck.ID, &deck.Name, &deck.EventID, &deck.UserID, &deck.LastScore)
		if err != nil {
			log.Fatal(err)
		}
		decks = append(decks, deck)
	}
	return decks, nil
}

func getCardsFromDB(memberID int64, deckID string) ([]*DB_Card, error) {
	query := `
		select deck_name, meetup_id, image_url from (
			select * from decks d
			inner join (
				select * from members
				where deck_id=$1
			) as m on m.deck_id=d.deck_id
		) as deck
		where user_id=$2
	`

	rows, err := db.Query(query, deckID, memberID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	cards := make([]*DB_Card, 0)
	for rows.Next() {
		card := new(DB_Card)
		err := rows.Scan(&card.Name, &card.MeetupID, &card.ImageURL)
		if err != nil {
			log.Fatal(err)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func addScoreToDB(memberID int64, deckID string, score string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM decks
		WHERE deck_id=$1 and user_id=$2
	`
	var count int

	err := db.QueryRow(query, deckID, memberID).Scan(&count)

	if err != nil {
		return false, err
	}

	if count != 1 {
		return false, nil
	}

	query = `
		INSERT INTO plays (deck_id, score)
		VALUES($1, $2)
		RETURNING deck_id
	`
	retID := 0
	err = db.QueryRow(query, deckID, score).Scan(&retID)

	if err != nil {
		return true, err
	}

	if retID == 0 {
		return false, nil
	}

	query = `
		UPDATE decks
		SET lastScore=$1
		WHERE deck_id=$2
		RETURNING deck_id
	`

	retID = 0
	err = db.QueryRow(query, score, deckID).Scan(&retID)

	if err != nil {
		return true, err
	}

	if retID == 0 {
		return false, nil
	}

	return true, nil
}

func setSelectedDeckInDB(memberID int64, deckID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM decks
		WHERE deck_id=$1 and user_id=$2
	`
	var count int

	err := db.QueryRow(query, deckID, memberID).Scan(&count)

	if err != nil {
		return false, err
	}

	if count != 1 {
		return false, nil
	}

	query = `
		INSERT INTO last_deck (deck_id, user_id)
		VALUES($1, $2)
		ON CONFLICT (user_id) DO
		UPDATE 
		SET deck_id=$1
		RETURNING user_id
	`

	retID := 0
	err = db.QueryRow(query, deckID, memberID).Scan(&retID)

	if err != nil {
		return false, err
	}
	if retID == 0 {
		return false, nil
	}

	return true, nil
}

func getLastDeckFromDB(memberID int64) (DB_Deck, error) {
	query := `
		select d.deck_id, d.deck_name, d.event_id, d.user_id, d.lastScore
		from decks d
		inner join (
			select * from last_deck
			where user_id=$1
		) as ld on d.deck_id=ld.deck_id
	`
	var deck DB_Deck

	err := db.QueryRow(query, memberID).Scan(&deck.ID, &deck.Name, &deck.EventID, &deck.UserID, &deck.LastScore)

	if err != nil {
		return deck, err
	}
	return deck, nil
}
