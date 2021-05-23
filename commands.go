package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// SetMyCommands this was used to setup the commands for the bot
func sendCommands(bot *Bot, commands []BotCommand) error {
	var msgToSend []BotCommand
	msgToSend = append(msgToSend, commands...)

	reqbytes, err := json.Marshal(msgToSend)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("commands", string(reqbytes))

	urlStr, err := makeURL(bot, SetMyCommands)
	if err != nil {
		return err
	}

	resp, err := makeRequest(bot, http.MethodPost, urlStr, data)
	if resp.StatusCode != http.StatusOK {
		return errors.New("StatusCode not Ok")
	}

	var hookResp WebhookResp
	if err := json.NewDecoder(resp.Body).Decode(&hookResp); err != nil {
		return err
	}

	if hookResp.Result == false || hookResp.OK == false {
		return errors.New("Webhook result not Ok")
	}

	return nil

}

// setupCommands this was used to setup the
func setupCommands(bot *Bot) error {
	var commands []BotCommand

	commands = []BotCommand{
		{
			Command:     "start",
			Description: "Initiates a conversation with telichess",
		},

		// {
		// 	Command:     "sayhello",
		// 	Description: "Greets you Hello",
		// },
		// {
		// 	Command:     "sendlocation",
		// 	Description: "Send my location to you",
		// },
		// {
		// 	Command:     "playgame",
		// 	Description: "We play a game to see if you are smarter than me.",
		// },
	}

	err := sendCommands(bot, commands)
	if err != nil {
		return err
	}
	return nil
}

//TODO: clean up this code to make it make sense
