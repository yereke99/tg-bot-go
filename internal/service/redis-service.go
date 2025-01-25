package service

import (
	"context"
	"tg-bot-go/config"

	"go.uber.org/zap"
)

type RedisService struct {
	ctx    context.Context
	config *config.Config
	logger *zap.Logger
}

func NewRedisService(ctx context.Context, logger *zap.Logger, config *config.Config) *RedisService {
	return &RedisService{
		ctx:    ctx,
		config: config,
		logger: logger,
	}
}
