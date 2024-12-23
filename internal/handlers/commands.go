package handlers

import (
	"asklyBot/internal/texts"
	"asklyBot/pkg/database"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"regexp"
	"strconv"
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
		Text:   texts.AddYourBirthday,
	})
}

func BirthdayInput(ctx context.Context, b *bot.Bot, update *models.Update) {
	userInput := update.Message.Text

	re := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	if !re.MatchString(userInput) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   texts.FailedToAddBirthday,
		})
		return
	}

	birthday := database.Birthday{
		ID:       strconv.FormatInt(update.Message.From.ID, 10),
		Name:     update.Message.From.FirstName,
		Birthday: userInput,
	}
	err := database.AddBirthday(birthday)
	if err != nil {
		log.Printf("Failed to save birthday: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to save your birthday",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   texts.Birthday,
	})
}
