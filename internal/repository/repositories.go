package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type IRedisRepository interface{}

type IChatRepository interface {
	AddUser(ctx context.Context, userID int64) error
	FindPartner(ctx context.Context, userID int64) (int64, error)
	SetPartner(ctx context.Context, userID, partnerID int64) error
	GetUserPartner(ctx context.Context, userID int64) (int64, error)
	RemoveUser(ctx context.Context, userID int64) error
	GetUsers(ctx context.Context) ([]int64, error)
}

type Repository struct {
	ChatRepository IChatRepository
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		ChatRepository: NewChatRepository(client),
	}
}
