package handler

import (
	"context"
	"fmt"
	"sync"
	"tg-bot-go/internal/keyboard"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ChatState struct {
	mu       sync.RWMutex
	users    []int64
	userChat map[int64]int64
}

var globalChatState *ChatState
var once sync.Once

func GetChatState() *ChatState {
	once.Do(func() {
		globalChatState = &ChatState{
			users:    make([]int64, 0),
			userChat: make(map[int64]int64),
		}
	})
	return globalChatState
}

func NewChat() *ChatState {
	return &ChatState{
		users:    make([]int64, 0),
		userChat: make(map[int64]int64),
	}
}

func (c *ChatState) AddUser(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, u := range c.users {
		if u == id {
			return // Пользователь уже в списке
		}
	}
	c.users = append(c.users, id)
}

func (c *ChatState) FindPartner(userID int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, user := range c.users {
		if user != userID {
			// Удаляем найденного партнёра из очереди
			c.users = append(c.users[:i], c.users[i+1:]...)
			return user
		}
	}
	return 0 // Если партнёров нет
}

func (c *ChatState) GetUsers() []int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.users
}

func (c *ChatState) SetPartner(id int64, partnerID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.userChat[id] = partnerID
}

func (c *ChatState) CheckPartnerToEmpty(id int64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.userChat[id]; ok {
		return true
	}
	return false
}

func (c *ChatState) GetUserPartner(id int64) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.userChat[id]
}

func (c *ChatState) RemoveUser(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Удаляем из списка users
	for i, u := range c.users {
		if u == id {
			c.users = append(c.users[:i], c.users[i+1:]...)
			break
		}
	}
	// Удаляем из карты пар
	delete(c.userChat, id)
}

func HandleChat(ctx context.Context, b *bot.Bot, update *models.Update, chatState *ChatState) {
	userID := update.Message.From.ID
	partnerID := chatState.GetUserPartner(userID)

	// Логируем входящее сообщение (общая информация)
	fmt.Printf("[LOG] UserID=%d -> PartnerID=%d | MessageType=", userID, partnerID)

	kbChat := keyboard.NewKeyboard()
	kbChat.AddRow(keyboard.NewInlineButton("💬 Chat", "chat"))

	if partnerID == 0 {
		// У пользователя нет собеседника
		fmt.Printf("NO_PARTNER\n") // логируем отсутствие партнёра
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:         update.Message.Chat.ID,
			Text:           "Вы не подключены к собеседнику. Нажмите кнопку 💬 Chat, чтобы начать.",
			ReplyMarkup:    kbChat.Build(),
			ProtectContent: true, // Делаем это сообщение приватным
		})
		return
	}

	// Кнопки для выхода из чата
	kb := keyboard.NewKeyboard()
	kb.AddRow(keyboard.NewInlineButton("🔕 Exit", "exit"))

	// Получаем username или FirstName
	username := update.Message.From.Username
	if username == "" {
		username = update.Message.From.FirstName
	}

	var caption string
	if update.Message.Caption != "" {
		caption = fmt.Sprintf("@%s: %s", username, update.Message.Caption)
	}

	// В зависимости от типа сообщения — логируем и пересылаем
	switch {
	case update.Message.Text != "":
		fmt.Printf("TEXT | User=@%s | Text=%q\n", username, update.Message.Text)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:         partnerID,
			Text:           fmt.Sprintf("@%s: %s", username, update.Message.Text),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Photo != nil:
		fmt.Printf("PHOTO | User=@%s | FileID=%s | Caption=%q\n",
			username,
			update.Message.Photo[len(update.Message.Photo)-1].FileID,
			update.Message.Caption,
		)

		photoID := update.Message.Photo[len(update.Message.Photo)-1].FileID
		b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:         partnerID,
			Photo:          &models.InputFileString{Data: photoID},
			Caption:        withDefaultCaption(username, caption, "фото"),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Video != nil:
		fmt.Printf("VIDEO | User=@%s | FileID=%s | Caption=%q\n",
			username,
			update.Message.Video.FileID,
			update.Message.Caption,
		)

		b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:         partnerID,
			Video:          &models.InputFileString{Data: update.Message.Video.FileID},
			Caption:        withDefaultCaption(username, caption, "видео"),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Voice != nil:
		fmt.Printf("VOICE | User=@%s | FileID=%s | Caption=%q\n",
			username,
			update.Message.Voice.FileID,
			update.Message.Caption,
		)

		b.SendVoice(ctx, &bot.SendVoiceParams{
			ChatID:         partnerID,
			Voice:          &models.InputFileString{Data: update.Message.Voice.FileID},
			Caption:        withDefaultCaption(username, caption, "голосовое сообщение"),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.VideoNote != nil:
		fmt.Printf("VIDEO_NOTE | User=@%s | FileID=%s\n",
			username,
			update.Message.VideoNote.FileID,
		)

		b.SendVideoNote(ctx, &bot.SendVideoNoteParams{
			ChatID:         partnerID,
			VideoNote:      &models.InputFileString{Data: update.Message.VideoNote.FileID},
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Document != nil:
		fmt.Printf("DOCUMENT | User=@%s | FileID=%s | Caption=%q\n",
			username,
			update.Message.Document.FileID,
			update.Message.Caption,
		)

		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:         partnerID,
			Document:       &models.InputFileString{Data: update.Message.Document.FileID},
			Caption:        withDefaultCaption(username, caption, "документ"),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Audio != nil:
		fmt.Printf("AUDIO | User=@%s | FileID=%s | Caption=%q\n",
			username,
			update.Message.Audio.FileID,
			update.Message.Caption,
		)

		b.SendAudio(ctx, &bot.SendAudioParams{
			ChatID:         partnerID,
			Audio:          &models.InputFileString{Data: update.Message.Audio.FileID},
			Caption:        withDefaultCaption(username, caption, "аудио"),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Location != nil:
		fmt.Printf("LOCATION | User=@%s | Lat=%.5f | Long=%.5f\n",
			username,
			update.Message.Location.Latitude,
			update.Message.Location.Longitude,
		)

		b.SendLocation(ctx, &bot.SendLocationParams{
			ChatID:         partnerID,
			Latitude:       update.Message.Location.Latitude,
			Longitude:      update.Message.Location.Longitude,
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Sticker != nil:
		fmt.Printf("STICKER | User=@%s | FileID=%s\n",
			username,
			update.Message.Sticker.FileID,
		)

		b.SendSticker(ctx, &bot.SendStickerParams{
			ChatID:         partnerID,
			Sticker:        &models.InputFileString{Data: update.Message.Sticker.FileID},
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Contact != nil:
		contact := update.Message.Contact
		fmt.Printf("CONTACT | User=@%s | Phone=%s | FirstName=%s | LastName=%s\n",
			username,
			contact.PhoneNumber,
			contact.FirstName,
			contact.LastName,
		)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: partnerID,
			Text: fmt.Sprintf("@%s отправил(а) контакт:\nТел: %s\nИмя: %s %s",
				username,
				contact.PhoneNumber,
				contact.FirstName,
				contact.LastName,
			),
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})

	case update.Message.Poll != nil:
		// Опрос
		poll := update.Message.Poll
		fmt.Printf("POLL | User=@%s | Question=%q | Options=%d\n",
			username,
			poll.Question,
			len(poll.Options),
		)
		// Переформируем варианты в InputPollOption
		var pollOptions []models.InputPollOption
		for _, o := range poll.Options {
			pollOptions = append(pollOptions, models.InputPollOption{Text: o.Text})
		}
		// Создаём новый опрос у собеседника
		b.SendPoll(ctx, &bot.SendPollParams{
			ChatID:         partnerID,
			Question:       poll.Question,
			Options:        pollOptions,
			ProtectContent: true,
			// Если хотите, можно добавлять: IsAnonymous: false и т.д.
		})

	default:
		fmt.Printf("UNKNOWN | User=@%s\n", username)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:         userID,
			Text:           "Неизвестный тип сообщения. Попробуйте отправить текст, фото, видео, голосовое сообщение или документ.",
			ReplyMarkup:    kb.Build(),
			ProtectContent: true,
		})
	}
}

// withDefaultCaption — вспомогательная функция:
// если у пользователя нет подписи (caption == ""),
// мы формируем её автоматически. Если есть, используем её.
func withDefaultCaption(username, caption, mediaType string) string {
	if caption != "" {
		return caption // Уже содержит @username и собственный текст
	}
	// Если подписи не было, можно составить свою
	return fmt.Sprintf("@%s отправил(а) %s", username, mediaType)
}
