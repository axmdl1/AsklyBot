package filters

import (
	"github.com/go-telegram/bot/models"
	"regexp"
)

const (
	start       = "/start"
	help        = "/help"
	addBirthday = "/addBirthday"
	api         = "/secretFunc"
)

func IsStart(update *models.Update) bool {
	return update.Message != nil && update.Message.Text == start
}

func IsHelp(update *models.Update) bool {
	return update.Message != nil && update.Message.Text == help
}

func IsBirthday(update *models.Update) bool {
	return update.Message != nil && update.Message.Text == addBirthday
}

func IsBirthdayInput(update *models.Update) bool {
	re := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	return update.Message != nil && re.MatchString(update.Message.Text)
}

func IsAPI(update *models.Update) bool {
	return update.Message != nil && update.Message.Text == api
}
