package main

import (
	"ASKLYBOT/clients/telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

}

func mustToken() string {
	token := flag.String(
		"token-bot",
		"",
		"Please give me token for Telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token not find")
	}

	return *token

}
