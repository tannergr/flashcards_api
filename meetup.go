package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

// RedirectURI is the meetup api redierect
const RedirectURI = "http://localhost:8080/api/meetup/auth"

// MeetupBase is the base url of the meetup auth flow
const MeetupBase = "https://secure.meetup.com/oauth2/"

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}

type jwtMeetup struct {
	AccessToken string
	jwt.StandardClaims
}

func buildMeetupToken(accessToken string) http.Cookie {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwtMeetup{
		AccessToken: accessToken,
	})
	tokenString, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		log.Fatal(err)
	}
	cookie := http.Cookie{Name: "jwt", Value: tokenString, Path: "/"}
	return cookie
}

func parseMeetupToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	accessTokenString, ok := claims[string("AccessToken")].(string)
	if ok {
		return accessTokenString, nil
	}
	return "not string", nil

}

func getMembersEvents(tokenString string) string {
	url := fmt.Sprintf("https://api.meetup.com/self/events?access_token=%v", tokenString)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		var events []Event
		err = json.Unmarshal([]byte(contents), &events)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(events); i++ {
			fmt.Printf("%s \n", events[i].Name)
		}
		jsonEvents, err := json.Marshal(events)
		if err != nil {
			log.Fatal(err)
		}
		return string(jsonEvents)
	}
	return "ok"
}

func getMemberInfo(tokenString string) Member {
	url := fmt.Sprintf("https://api.meetup.com/members/self?access_token=%v", tokenString)
	response, err := http.Get(url)
	var member Member
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		// for k, v := range response.Header {
		// 	fmt.Print(k)
		// 	fmt.Print(" : ")
		// 	fmt.Println(v)
		// }
		contents, err := ioutil.ReadAll(response.Body)
		// fmt.Printf(string(contents))
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		err = json.Unmarshal([]byte(contents), &member)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name: %s \n", member.Name)
		fmt.Printf("MemberID: %v \n", member.ID)
		return member
	}
	return member
}

func getMeetupEventMembers(tokenString string, groupname string, eid string) string {
	url := fmt.Sprintf("https://api.meetup.com/%v/events/%v/attendance?access_token=%v",
		groupname,
		eid,
		tokenString)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		return string(contents)
	}
	return "ok"
}

func getMeetupEventMembersArray(tokenString string, groupname string, eid string) []Card {
	var RSVPs []MeetupRSVP
	var cards []Card
	url := fmt.Sprintf("https://api.meetup.com/%v/events/%v/attendance?access_token=%v",
		groupname,
		eid,
		tokenString)
	fmt.Println(url)

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		err = json.Unmarshal([]byte(contents), &RSVPs)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(RSVPs); i++ {
			//name, imageurl, id
			name := RSVPs[i].Member.Name
			image := RSVPs[i].Member.Photo.ImageURL
			id := RSVPs[i].Member.MeetupID
			cards = append(cards, Card{name, image, id})
		}
		return cards
	}
	return cards
}
