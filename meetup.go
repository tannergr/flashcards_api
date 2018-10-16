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

func parseMeetupToken(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	accessTokenString, ok := claims[string("AccessToken")].(string)
	if ok {
		return accessTokenString
	}
	return "not string"

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
		contents, err := ioutil.ReadAll(response.Body)
		fmt.Printf(string(contents))
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		err = json.Unmarshal([]byte(contents), &member)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s \n", member.Name)
		fmt.Printf("%v \n", member.ID)
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
