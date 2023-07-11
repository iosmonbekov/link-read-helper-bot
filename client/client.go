package client

type OkResponse struct {
	Ok bool `json:"ok"`
}

type ResultResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type User struct {
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}

type Client interface {
	FetchUpdates(offset int, limit int) ([]Update, error)
	SendMessage(update Update, message string) error
}
