package slack

import "encoding/json"

type Message struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func (m Message) ToJSON() string {
	buf, _ := json.Marshal(&m)
	return string(buf)
}

func NewMessage(text string) Message {
	return Message{
		Text:         text,
		ResponseType: "in_channel",
	}
}
