package main

type BaseMessage struct {
	Type string `json:"type"`
}

type Register struct {
	Type   string `json:"type"`
	UserID int    `json:"user_id"`
}

type SubMain struct {
	Type         string `json:"type"`
	Subscription bool   `json:"subscription"`
}

type UpdateCountUser struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}
