package main

// для визначення типу повідомлення
type BaseMessage struct {
	Type string `json:"type"`
}

//надіслати користувачу його айді

type SubMain struct {
	Type         string `json:"type"`
	Subscription bool   `json:"subscription"`
}

type UpdateCountUser struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type FindInterlocutor struct {
	Type string `json:"type"`
}

type TextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type DeleteNotice struct {
	Type string `json:"type"`
}
