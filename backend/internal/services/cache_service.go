package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheService(redisURL string) *CacheService {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil
	}

	return &CacheService{
		client: client,
		ctx:    ctx,
	}
}

func (s *CacheService) Get(key string, dest interface{}) error {
	val, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (s *CacheService) Set(key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(s.ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

func (s *CacheService) Delete(key string) error {
	return s.client.Del(s.ctx, key).Err()
}

func (s *CacheService) DeletePattern(pattern string) error {
	iter := s.client.Scan(s.ctx, 0, pattern, 0).Iterator()
	for iter.Next(s.ctx) {
		if err := s.client.Del(s.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (s *CacheService) Exists(key string) (bool, error) {
	result, err := s.client.Exists(s.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (s *CacheService) GetStats() (map[string]interface{}, error) {
	info, err := s.client.Info(s.ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	dbSize, err := s.client.DBSize(s.ctx).Result()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"info":    info,
		"db_size": dbSize,
	}, nil
}

func (s *CacheService) FlushAll() error {
	return s.client.FlushAll(s.ctx).Err()
}

func (s *CacheService) InvalidateBooks() error {
	return s.DeletePattern("books:*")
}

func (s *CacheService) InvalidateCategories() error {
	return s.DeletePattern("categories:*")
}

func (s *CacheService) InvalidateAuthors() error {
	return s.DeletePattern("authors:*")
}

func (s *CacheService) InvalidateBlogs() error {
	return s.DeletePattern("blogs:*")
}

func (s *CacheService) InvalidateFAQs() error {
	return s.DeletePattern("faqs:*")
}

func (s *CacheService) InvalidateTestimonials() error {
	return s.DeletePattern("testimonials:*")
}

func (s *CacheService) InvalidateSettings() error {
	return s.DeletePattern("settings:*")
}

func (s *CacheService) InvalidateUserPermissions(userID uint) error {
	return s.Delete("user:permissions:" + string(rune(userID)))
}
