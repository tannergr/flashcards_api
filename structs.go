package main

type Deck struct {
	ID        int64  `json:"deck_id"`
	Name      string `json:"name"`
	EventID   string `json:"event_id"`
	Cards     []Card `json:"cards"`
	GroupName string `json:"group_name"`
}

type DB_Deck struct {
	ID        int64  `json:"deck_id"`
	Name      string `json:"name"`
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	LastScore int    `json:"lastScore"`
	GroupName string `json:"group_name"`
}

type Card struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageurl"`
	MeetupID int    `json:"meetup_id"`
}

type GameRound struct {
	Names  []string `json:"names"`
	Actual int      `json:"actual"`
	Photo  string   `json:"photo"`
}

type DB_Card struct {
	MeetupID int `json:"meetup_id"`
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

type MeetupRSVP struct {
	Member Event_Member `json:"member"`
}

type Event_Member struct {
	Name     string `json:"name"`
	Photo    Photo  `json:"photo"`
	MeetupID int    `json:"id"`
}

type Photo struct {
	ImageURL string `json:"highres_link"`
}
