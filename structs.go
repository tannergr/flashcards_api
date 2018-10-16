package main

type Deck struct {
	ID      int64  `json:"deck_id"`
	Name    string `json:"name"`
	EventID string `json:"event_id"`
	Cards   []Card `json:"cards"`
}

type DB_Deck struct {
	ID      int64  `json:"deck_id"`
	Name    string `json:"name"`
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
}

type Card struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageurl"`
}

type DB_Card struct {
	Name     string `json:"deck_name"`
	ImageURL string `json:"imageurl"`
	MeetupID string `json:"meetup_id"`
}

type Member struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	URLNAME string `json:"urlname"`
}

type Event struct {
	LocalTime string `json:"local_time"`
	LocalDate string `json:"local_date"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Group     Group  `json:"group"`
}

type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"`
}

type expirationTime int32

type errResp struct {
	Error string
}

type succResp struct {
	Message string
}
