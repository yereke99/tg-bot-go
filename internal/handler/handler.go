package handler

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StateID string

const (
	StateDefault StateID = "default"
	StateAskName StateID = "ask_name"
	StateAskAge  StateID = "ask_age"
	StateFinish  StateID = "finish"
)

type User struct {
	State, Name StateID
	Age         int
}

type UserState struct {
	mu    sync.RWMutex
	users map[int64]*User
}

func NewUserState() *UserState {
	return &UserState{
		users: make(map[int64]*User),
	}
}

func (us *UserState) GetUser(id int64) *User {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, ok := us.users[id]
	if !ok {
		return nil
	}
	return user
}

func (us *UserState) SetUser(id int64, user *User) {
	us.mu.Lock()
	defer us.mu.Unlock()

	us.users[id] = user
}

func (us *UserState) GetState(id int64) StateID {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, ok := us.users[id]
	if !ok {
		return StateDefault
	}
	return user.State
}

func (us *UserState) SetState(id int64, newState StateID) {
	us.mu.Lock()
	defer us.mu.Unlock()

	_, ok := us.users[id]
	if !ok {
		user := &User{}
		us.users[id] = user
	}
	us.users[id].State = newState
}

func HandleUpdate(ctx context.Context, b *bot.Bot, update *models.Update, us *UserState) {

	if update.Message == nil || update.Message.Text == "" {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	message := update.Message.Text

	switch message {
	case "/form":
		us.SetState(userID, StateAskName)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      "Please enter your name:",
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}

	state := us.GetState(userID)
	handleStateTransition(ctx, b, update, us, state)

}

func handleStateTransition(ctx context.Context, b *bot.Bot, update *models.Update, us *UserState, state StateID) {
	switch state {
	case StateAskName:
		handleAskName(ctx, b, update, us)
	case StateAskAge:
		handleAskAge(ctx, b, update, us)
	case StateFinish:
		handleFinishedState(ctx, b, update, us)
	default:
		handleDefaultState(ctx, b, update)
	}
}

func handleDefaultState(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Please start the form by typing /form.",
	})
}

func handleAskName(ctx context.Context, b *bot.Bot, update *models.Update, us *UserState) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	name := update.Message.Text

	if len(name) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Invalid name. Please enter a name of at least 2 characters.",
		})
		return
	}

	user := us.GetUser(userID)
	if user == nil {
		user = &User{State: StateDefault}
	}
	user.Name = StateID(name)

	us.SetUser(userID, user)
	us.SetState(userID, StateAskAge)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Great! Now please enter your age.",
	})
}

func handleAskAge(ctx context.Context, b *bot.Bot, update *models.Update, us *UserState) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	ageInput := update.Message.Text

	age, err := strconv.Atoi(ageInput)
	if err != nil || age < 18 || age > 100 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Please enter a valid age between 18 and 100.",
		})
		return
	}

	user := us.GetUser(userID)
	if user == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Something went wrong. Please type /form to start over.",
		})
		us.SetState(userID, StateDefault)
		return
	}

	user.Age = age
	us.SetUser(userID, user)
	us.SetState(userID, StateFinish)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Thank you! Your name is %s and your age is %d.", user.Name, user.Age),
	})
}

func handleFinishedState(ctx context.Context, b *bot.Bot, update *models.Update, us *UserState) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	us.SetState(userID, StateDefault)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "You have already completed the form. Type /form to fill it out again.",
	})
}
