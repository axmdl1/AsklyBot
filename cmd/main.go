package main

import (
	"context"
	"os"
	"os/signal"
	"tgBotExample/internal/filters"
	"tgBotExample/internal/handlers"
	"tgBotExample/pkg/systems"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started

func main() {
	token := systems.BotToken()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandlerMatchFunc(filters.IsStart, handlers.Start)
	b.RegisterHandlerMatchFunc(filters.IsHelp, handlers.Help)
	b.RegisterHandlerMatchFunc(filters.IsBirthday, handlers.AddBirthday)

	b.RegisterHandlerMatchFunc(filters.IsPhoto, handlers.Photo)
	b.RegisterHandlerMatchFunc(filters.IsVideo, handlers.Video)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
