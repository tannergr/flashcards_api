package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Handles call back from meetup login
// retireves the access code from meetup
// returns a jwt and redirects
var CallBackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("code:", getParam(r, "code"))
	reqURL := fmt.Sprintf(`%vaccess?`+
		`grant_type=authorization_code`+
		`&code=%v`+
		`&redirect_uri=%v`+
		`&client_id=%v`+
		`&client_secret=%v`,
		MeetupBase, getParam(r, "code"), RedirectURI,
		ClientID, ClientSecret)
	fmt.Println(reqURL)
	resp, err := http.Post(reqURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		fmt.Printf("HTTP request failed")
	} else {
		defer resp.Body.Close()
		var accessToken tokenJSON
		err := json.NewDecoder(resp.Body).Decode(&accessToken)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Access Token:", accessToken.AccessToken)

		cookie := buildMeetupToken(accessToken.AccessToken)
		http.SetCookie(w, &cookie)
	}

	fmt.Fprintf(w, formatRequest(r))
})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href=%vauthorize?`+
		`&client_id=%v`+
		`&redirect_uri=%v`+
		`&response_type=code>click me</a>`,
		MeetupBase, ClientID, RedirectURI)
})

var GetEvents = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]

	if cookie == "" {
		sendUnathorized(w)
		return
	}
	verifiedToken := parseMeetupToken(cookie)

	events := getMembersEvents(verifiedToken)
	fmt.Fprintf(w, events)
})

var getEventMembers = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]

	if cookie == "" {
		sendUnathorized(w)
		return
	}

	verifiedToken := parseMeetupToken(cookie)
	vars := mux.Vars(r)

	members := getMeetupEventMembers(verifiedToken, vars["groupname"], vars["eid"])
	fmt.Fprintf(w, members)
})

var GetMember = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookieName := "jwt"
	cookie, err := r.Cookie(cookieName)

	if err != nil {
		panic(err)
	}
	fmt.Printf(cookie.Value)
	verifiedToken := parseMeetupToken(cookie.Value)

	member := getMemberInfo(verifiedToken)
	jsonMember, err := json.Marshal(member)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(jsonMember))
	return
})

var AddDeck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]

	if cookie == "" {
		sendUnathorized(w)
		return
	}
	verifiedToken := parseMeetupToken(cookie)

	member := getMemberInfo(verifiedToken)
	decoder := json.NewDecoder(r.Body)
	var deck Deck
	err := decoder.Decode(&deck)
	fmt.Println(formatRequest(r))
	if err != nil {
		panic(err)
	}

	// Create new Deck in deck database
	// - memberID
	// - deck Name

	// Create new cards with deck id
	// - deck ID
	// - card member
	// - imageURL
	err = addDeckToDB(member.ID, deck)

	if err != nil {
		panic(err)
	}
	sendMessage(w, "success")
	return
})

var GetDecks = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]

	verifiedToken := parseMeetupToken(cookie)

	member := getMemberInfo(verifiedToken)

	decks, err := getDecksFromDB(member.ID)
	if err != nil {
		panic(err)
	}

	jsonDecks, err := json.Marshal(decks)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(jsonDecks))
	return
})

var GetCards = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]
	verifiedToken := parseMeetupToken(cookie)
	member := getMemberInfo(verifiedToken)
	vars := mux.Vars(r)
	cards, err := getCardsFromDB(member.ID, vars["deckID"])

	if err != nil {
		panic(err)
	}

	jsonCards, err := json.Marshal(cards)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(jsonCards))

	return
})

var PostScore = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]
	verifiedToken := parseMeetupToken(cookie)
	member := getMemberInfo(verifiedToken)
	vars := mux.Vars(r)

	match, err := addScoreToDB(member.ID, vars["deckID"], vars["score"])

	if err != nil {
		panic(err)
	}
	if match == false {
		fmt.Fprintf(w, "{message:'no matching deck for user'}")
		return
	}

	fmt.Fprintf(w, "{message:'success'}")

	return
})

var SetSelectedDeck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]
	verifiedToken := parseMeetupToken(cookie)
	member := getMemberInfo(verifiedToken)
	vars := mux.Vars(r)

	match, err := setSelectedDeckInDB(member.ID, vars["deckID"])

	if err != nil {
		panic(err)
	}
	if match == false {
		fmt.Fprintf(w, "{message:'no matching deck for user'}")
		return
	}

	fmt.Fprintf(w, "{message:'success'}")

	return
})

var GetLastDeck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header["Authentication"][0]
	verifiedToken := parseMeetupToken(cookie)
	member := getMemberInfo(verifiedToken)
	fmt.Println("getting last deck")
	deck, err := getLastDeckFromDB(member.ID)
	if err != nil {
		panic(err)
	}
	jsonDeck, err := json.Marshal(deck)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(jsonDeck))

	return
})
