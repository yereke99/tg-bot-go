package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	StateDefaults StateID = "default"
	StateChat     StateID = "chat"
	StateExit     StateID = "exit"
)

type Chat struct {
	userChat  map[int64]int64
	userState map[int64]StateID
}

type ChatState struct {
	mu    sync.RWMutex
	users []int64
	room  *Chat
}

func NewChat() *ChatState {
	room := make(map[int64]int64)
	userState := make(map[int64]StateID)
	chat := &Chat{
		userChat:  room,
		userState: userState,
	}
	return &ChatState{
		users: make([]int64, 0, 1000),
		room:  chat,
	}
}

func (c *ChatState) GetUsers() []int64 {
	return c.users
}

func (c *ChatState) AddUser(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.check(id) {
		fmt.Println("user is already exist")
		return
	}
	c.users = append(c.users, id)
}

func (c *ChatState) check(id int64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := 0; i < len(c.users); i++ {
		if id == c.users[i] {
			return true
		}
	}
	return false
}

func (c *ChatState) GetState(id int64) StateID {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if state, ok := c.room.userState[id]; ok {
		return state
	}

	return StateDefaults
}

func (c *ChatState) GetUserPartner(id int64) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if partner, ok := c.room.userChat[id]; ok {
		return partner
	}
	return 0
}

func (c *ChatState) SetState(id int64, status StateID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.room.userState[id] = status
}

func (c *ChatState) SetPartner(id int64, partnerId int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.room.userChat[id] = partnerId
}

func HandleChat(ctx context.Context, b *bot.Bot, update *models.Update, userChat *ChatState) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	message := update.Message.Text

	switch message {
	case "/chat":
		_, _ = userID, chatID
	}
}
