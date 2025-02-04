package repository

import (
	"context"
	"testing"
	"tg-bot-go/config"

	"github.com/stretchr/testify/assert"
)

func setupTestRedisRepo() *RedisRepository {
	cfg := &config.Config{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       1, // Используем отдельную тестовую базу
	}
	return NewRedisRepository(cfg)
}

func TestRedisRepository_SetAndGetUserState(t *testing.T) {
	ctx := context.Background()
	redisRepo := setupTestRedisRepo()
	userID := int64(12345)
	expectedState := "ask_name"

	// Тестируем SetUserState
	err := redisRepo.SetUserState(ctx, userID, expectedState)
	assert.NoError(t, err, "SetUserState должен завершаться без ошибок")

	// Тестируем GetUserState
	actualState, err := redisRepo.GetUserState(ctx, userID)
	assert.NoError(t, err, "GetUserState должен завершаться без ошибок")
	assert.Equal(t, expectedState, actualState, "Состояние пользователя должно совпадать")
}

func TestRedisRepository_SetAndGetUserDatas(t *testing.T) {
	ctx := context.Background()
	redisRepo := setupTestRedisRepo()
	userID := int64(12345)

	expectedData := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	err := redisRepo.SetUserData(ctx, userID, expectedData)
	assert.NoError(t, err, "SetUserData должен завершаться без ошибок")

	actualData, err := redisRepo.GetUserData(ctx, userID)
	assert.NoError(t, err, "GetUserData должен завершаться без ошибок")
	assert.Equal(t, "John", actualData["name"], "Имя должно совпадать")
	assert.Equal(t, "30", actualData["age"], "Возраст должен совпадать")
}

func TestRedisRepository_GetUserState_NonExistentKey(t *testing.T) {
	ctx := context.Background()
	redisRepo := setupTestRedisRepo()
	userID := int64(99999) // Несуществующий ключ

	// Тестируем GetUserState для несуществующего ключа
	state, err := redisRepo.GetUserState(ctx, userID)
	assert.NoError(t, err, "GetUserState для несуществующего ключа должен завершаться без ошибок")
	assert.Equal(t, "", state, "Состояние должно быть пустой строкой для несуществующего ключа")
}

func TestRedisRepository_GetUserData_NonExistentKey(t *testing.T) {
	ctx := context.Background()
	redisRepo := setupTestRedisRepo()
	userID := int64(99999) // Несуществующий ключ

	// Тестируем GetUserData для несуществующего ключа
	data, err := redisRepo.GetUserData(ctx, userID)
	assert.NoError(t, err, "GetUserData для несуществующего ключа должен завершаться без ошибок")
	assert.Empty(t, data, "Данные должны быть пустыми для несуществующего ключа")
}
