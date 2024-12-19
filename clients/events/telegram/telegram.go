package telegram

import "ASKLYBOT/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
}
