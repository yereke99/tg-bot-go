package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// HandlePollAnswer - вызывается, когда кто-то проголосовал в опросе
func HandlePollAnswer(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Проверяем, что это PollAnswer
	if update.PollAnswer == nil {
		return
	}

	answer := update.PollAnswer
	userID := answer.User.ID
	pollID := answer.PollID
	chosenOptions := answer.OptionIDs

	// Например, отправим этому же пользователю инфо о его выборе
	msg := fmt.Sprintf("Вы проголосовали в опросе ID=%s. Варианты: %v", pollID, chosenOptions)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   msg,
	})

	fmt.Printf("[LOG] User=%d выбрал варианты %v в опросе %s\n",
		userID,
		chosenOptions,
		pollID,
	)
}
