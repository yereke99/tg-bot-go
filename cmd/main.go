package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"tg-bot-go/internal/handler"
	"tg-bot-go/internal/keyboard"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Replace with your bot token
	token := "1325617758:AAHD8tkdxsDOE2M5oAP9BW5LF71dg5KdRQo"

	opts := []bot.Option{
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackHandler),
		bot.WithCallbackQueryDataHandler("chat", bot.MatchTypePrefix, callbackHandler),
	}

	// Create bot
	b, err := bot.New(token, opts...)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	// Create user state manager
	//userState := handler.NewUserState()

	// Register a single text message handler
	//b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *models.Update) {
	//	handler.HandleUpdate(ctx, b, update, userState)
	//})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/hello", bot.MatchTypeExact, helloHandler)

	fmt.Println("Bot is running...")
	b.Start(ctx)
}

func helloHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := keyboard.NewKeyboard()
	kb.AddRow(keyboard.NewInlineButton("💬 Chat", "chat"))

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Hello, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + " please click button 💬 Chat to join room*",
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: kb.Build(),
	})
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       true,
	})
	userID := update.Message.From.ID

	chat := handler.NewChat()
	chat.AddUser(userID)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}
