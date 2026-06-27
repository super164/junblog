package database

import (
	"context"
	"fmt"

	"blog_backend/pkg/config"
	"blog_backend/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var redisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	// 如果没有配置 Redis，返回 nil
	if cfg.Host == "" {
		logger.Warn("Redis 未配置，跳过初始化")
		return nil, nil
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Redis 连接失败", zap.Error(err))
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	logger.Info("Redis 连接成功",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return redisClient, nil
}

// GetRedis 获取 Redis 实例
func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
