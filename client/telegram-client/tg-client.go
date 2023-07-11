package telegram_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"link-saver-bot/client"
	"net/http"
	"net/url"
	"strconv"
)

type TelegramClient struct {
	host     string
	basePath string
	token    string
}

func New(host string, token string) TelegramClient {
	return TelegramClient{
		host:     host,
		basePath: basePath(token),
		token:    token,
	}
}

func (c TelegramClient) FetchUpdates(offset int, limit int) ([]client.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	resp, err := http.Get(c.host + c.basePath + "/getUpdates?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("telegram_client.FetchUpdates: error on http.Get -> %v\n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("telegram_client.FetchUpdates: error on ioutil.ReadAll -> %v\n", err)
	}

	var response client.ResultResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("telegram_client.FetchUpdates: error on json.Unmarshal -> %v\n", err)
	}

	return response.Result, nil
}

func (c TelegramClient) SendMessage(update client.Update, message string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(update.Message.Chat.Id))
	q.Add("text", message)

	resp, err := http.Get(c.host + c.basePath + "/sendMessage?" + q.Encode())
	if err != nil {
		return fmt.Errorf("telegram_client.SendMessage: error on http.Get -> %v\n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("telegram_client.SendMessage: error on ioutil.ReadAll -> %v\n", err)
	}

	var response client.OkResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("telegram_client.SendMessage: error on json.Unmarshal -> %v\n", err)
	}

	if !response.Ok {
		return fmt.Errorf("telegram_client.SendMessage: response.Ok is false\n")
	}

	return nil
}

func basePath(token string) string {
	return "/bot" + token
}
