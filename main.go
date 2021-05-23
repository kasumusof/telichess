package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//wallstreetbet

var (
	err    error
	update *Update
	port   string
	token  string
)

func init() {
	colors[0] = "White"
	colors[1] = "Black"
	log.Println("Starting telichess...")
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error initislizing the local environment", err)
	}
	token = os.Getenv("TOKEN")
	if token == "" {
		log.Panic("Set TOKEN in your environment variables")
	}
	port = os.Getenv("PORT")
	if token == "" {
		log.Panic("Set PORT in your environment variables")
	}
	myBot, err = NewBot(token)
	if err != nil {
		log.Panic("Error initilalizing bot:", err)
	}
}

func main() {
	log.Println("telichess started...")
	go handleUpdate(updateChan)
	http.HandleFunc(fmt.Sprintf("/guess/%s", token), webHookUpdate)
	if err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Panic("Error listening by http package:", err)
	}
	// log.Println("Bot is  running.")
}

func webHookUpdate(w http.ResponseWriter, r *http.Request) {
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Println("could not read update Body:", err)
	}
	updateChan <- update
}
