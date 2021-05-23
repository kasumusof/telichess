package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func sendMessage(bot *Bot, chatID int, text string, keys [][]InlineKeyboardButton) (*Message, error) {

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: keys,
	}

	msgToSend := &SendMessageBody{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	}

	data := structToMap(msgToSend)

	urlStr, err := makeURL(bot, SendMessage)
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(bot, http.MethodPost, urlStr, data)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode not Ok")
	}

	msgReceived := &Message{}
	err = json.NewDecoder(resp.Body).Decode(msgReceived)
	if err != nil {
		return nil, err
	}
	return msgReceived, nil
}

func editMessageText(bot *Bot, chatID, messageID int, inlineMessageID, text string, keys [][]InlineKeyboardButton) (*Message, error) {

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: keys,
	}
	msgToSend := &messageTextEdit{
		MessageID:   messageID,
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	}
	data := structToMap(msgToSend)
	// a, _ := json.Marshal(data)
	// fmt.Println("Data to be sent", string(a))

	urlStr, err := makeURL(bot, EditMessageText)
	if err != nil {
		return nil, err
	}
	// fmt.Println("URL", urlStr)

	resp, err := makeRequest(bot, http.MethodPost, urlStr, data)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode not Ok button")
	}

	msgReceived := &Message{}
	err = json.NewDecoder(resp.Body).Decode(msgReceived)
	if err != nil {
		return nil, err
	}
	// fmt.Println(msgReceived)
	return msgReceived, nil
}

func answerCallbackQuery(bot *Bot, chatID, text string, time int) (*Message, error) {
	msgToSend := &callBackQueryAnswer{
		CallbackQueryID: chatID,
		Text:            text,
		ShowAlert:       false,
		CacheTime:       time,
	}

	data := structToMap(msgToSend)
	// b, _ := json.Marshal(data)
	// fmt.Println("Data to be sent", string(b))

	urlStr, err := makeURL(bot, AnswerCallbackQuery)
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(bot, http.MethodPost, urlStr, data)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode not Ok button")
	}

	msgReceived := &Message{}
	err = json.NewDecoder(resp.Body).Decode(msgReceived)
	if err != nil {
		return nil, err
	}
	// fmt.Println(msgReceived)
	return msgReceived, nil
}
