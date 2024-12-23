package handlers

import (
	"asklyBot/internal/texts"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   texts.Welcome,
	})
}

// It's a temporary function, it will be deleted after working with api
// Make sure that user cannot get this function
func GettingAPI(ctx context.Context, b *bot.Bot, update *models.Update) {
	jsonMessage, err := json.MarshalIndent(update.Message, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   string(jsonMessage),
	})
	fmt.Println(update.Message)
}

func Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   texts.Help,
	})
}

func AddBirthday(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Enter date:",
	})
	//needs to asking for birthday
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   texts.Birthday,
	})
	//successfully message to user
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   texts.Birthday,
	})
}
