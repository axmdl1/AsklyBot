package main

import (
	"flag"
	"log"
)

func main() {
	t := mustToken()
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
