package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"tg-bot-go/config"
	"tg-bot-go/internal/handler"
	"tg-bot-go/internal/keyboard"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Replace with your bot token
	token := cfg.Token

	opts := []bot.Option{
		bot.WithCallbackQueryDataHandler("chat", bot.MatchTypePrefix, chatButtonHandler),
		bot.WithCallbackQueryDataHandler("select_", bot.MatchTypePrefix, inlineHandler),
		bot.WithCallbackQueryDataHandler("exit", bot.MatchTypePrefix, callbackHandlerExit),
	}

	// Create bot
	b, err := bot.New(token, opts...)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	chatState := handler.GetChatState()

	// 1) Регистрируем хендлер для обычных сообщений (пересылка между собеседниками)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *models.Update) {
			handler.HandleChat(ctx, b, update, chatState)
		})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/hello", bot.MatchTypeExact, helloHandler)

	// Create user state manager
	//userState := handler.NewUserState()

	// Register a single text message handler
	//b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, func(ctx context.Context, b *bot.Bot, update *models.Update) {
	//	handler.HandleUpdate(ctx, b, update, userState)
	//})

	fmt.Println("Bot is running...")
	b.Start(ctx)
}

func inlineHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatState := handler.GetChatState() // Получаем глобальное состояние чата
	userID := update.CallbackQuery.From.ID

	// Извлекаем ID выбранного пользователя
	var selectedUserID int64
	fmt.Sscanf(update.CallbackQuery.Data, "select_%d", &selectedUserID)

	if chatState.CheckPartnerToEmpty(selectedUserID) {
		// Уведомляем пользователей что пользователь занят!
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   fmt.Sprintf("Собеседник сейчас занят, пожалуйста подождите: %d", selectedUserID),
		})
		return
	}
	// Устанавливаем партнёров
	chatState.SetPartner(userID, selectedUserID)
	chatState.SetPartner(selectedUserID, userID)

	// Уведомляем пользователей
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   fmt.Sprintf("Вы подключены к собеседнику с ID: %d", selectedUserID),
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: selectedUserID,
		Text:   fmt.Sprintf("Вы подключены к собеседнику с ID: %d", userID),
	})
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

func callbackHandlerExit(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatState := handler.GetChatState()

	userID := update.CallbackQuery.From.ID
	partnerID := chatState.GetUserPartner(userID)

	kb := keyboard.NewKeyboard()
	kb.AddRow(keyboard.NewInlineButton("💬 Chat", "chat"))

	chatState.RemoveUser(userID)
	if partnerID != 0 {
		chatState.RemoveUser(partnerID)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      partnerID,
			Text:        "Ваш собеседник покинул чат.",
			ReplyMarkup: kb.Build(),
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.From.ID,
		Text:        "Вы вышли из чата.",
		ReplyMarkup: nil,
	})
}

func chatButtonHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatState := handler.GetChatState() // Получаем глобальное состояние чата
	userID := update.CallbackQuery.From.ID

	// Добавляем пользователя в список, если его там нет
	chatState.AddUser(userID)

	// Формируем список пользователей в виде инлайн-кнопок
	users := chatState.GetUsers()
	kb := keyboard.NewKeyboard()
	for _, u := range users {
		if u != userID { // Исключаем самого пользователя
			kb.AddRow(keyboard.NewInlineButton(fmt.Sprintf("User %d", u), fmt.Sprintf("select_%d", u)))
		}
	}

	// Если список
	if len(users) == 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Нет доступных пользователей для подключения. Подождите...",
		})
		return
	}

	// Отправляем сообщение с инлайн-кнопками
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.From.ID,
		Text:        "Выберите пользователя для подключения:",
		ReplyMarkup: kb.Build(),
	})
}
