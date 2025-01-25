package repository

import (
	"context"
	"strconv"
	"tg-bot-go/config"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository создаёт новый экземпляр RedisRepository с конфигурацией.
func NewRedisRepository(cfg *config.Config) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,     // Адрес Redis
		Password: cfg.RedisPassword, // Пароль Redis (если требуется)
		DB:       cfg.RedisDB,       // Номер базы данных Redis
	})
	return &RedisRepository{client: client}
}

// SetUserState сохраняет состояние пользователя в Redis.
func (r *RedisRepository) SetUserState(ctx context.Context, userID int64, state string) error {
	key := strconv.FormatInt(userID, 10) + "_state" // Добавляем суффикс _state
	return r.client.Set(ctx, key, state, 0).Err()
}

// GetUserState получает состояние пользователя из Redis.
func (r *RedisRepository) GetUserState(ctx context.Context, userID int64) (string, error) {
	key := strconv.FormatInt(userID, 10) + "_state" // Добавляем суффикс _state
	state, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Если ключ не найден, возвращаем пустую строку
	}
	return state, err
}

// SetUserData сохраняет данные пользователя в Redis.
func (r *RedisRepository) SetUserData(ctx context.Context, userID int64, data map[string]interface{}) error {
	key := strconv.FormatInt(userID, 10) + "_data" // Добавляем суффикс _data
	return r.client.HSet(ctx, key, data).Err()
}

// GetUserData получает данные пользователя из Redis.
func (r *RedisRepository) GetUserData(ctx context.Context, userID int64) (map[string]string, error) {
	key := strconv.FormatInt(userID, 10) + "_data" // Ключ - только userID с суффиксом _data
	data, err := r.client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Если ключ не найден, возвращаем nil
	}
	if len(data) == 0 {
		return nil, nil // Если данных нет, возвращаем nil
	}
	return data, err
}
