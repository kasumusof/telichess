package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	myBot *Bot
)

// Bot struct
// This represents a Bot object that'll make api calls to the telegram api
type Bot struct {
	Name     string
	Token    string
	Client   *http.Client
	APIURL   string
	endpoint string
}

// NewBot creates a new bot instance
func NewBot(token string) (*Bot, error) {
	bot := &Bot{
		Name:     "teliChess",
		Token:    token,
		Client:   &http.Client{},
		APIURL:   "https://api.telegram.org/",
		endpoint: "bot%s/",
	}
	bot.endpoint = fmt.Sprintf(bot.endpoint, bot.Token)
	if err := bot.ready(); err != nil {
		return bot, err
	}

	return bot, nil
}

// Ready method sets the webhook on telegram
func (bot *Bot) ready() error {
	URL := os.Getenv("WEBHOOKURL")
	if URL == "" {
		return errors.New("environment variable WEBHOOKURL not set")
	}
	URL = fmt.Sprintf("%s/guess/%s", URL, token)
	data := url.Values{}
	data.Set("url", URL)

	urlStr, err := makeURL(bot, SetWebhook)
	if err != nil {
		return err
	}

	resp, err := makeRequest(bot, http.MethodPost, urlStr, data)
	if err != nil {
		log.Println()
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("StatusCode not Ok")
	}

	var hookResp WebhookResp
	if err := json.NewDecoder(resp.Body).Decode(&hookResp); err != nil {
		return err
	}

	if hookResp.Result == false {
		return errors.New("Webhook result not Ok")
	}

	err = setupCommands(bot)
	if err != nil {
		return err
	}
	return nil
}
