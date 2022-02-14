package repository

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type RedisStorage interface {
	Get(key string) (interface{}, error)
	SaveOrCreate(key string, age int, value []byte) error
	Delete(key string) error
}

type redisStorage struct {
	redis   *redis.Client
	context context.Context
}

func NewRedisStorage(rc *redis.Client, ctx context.Context) RedisStorage {
	return &redisStorage{
		redis:   rc,
		context: ctx,
	}
}

func (re *redisStorage) Get(key string) (interface{}, error) {
	v, err := re.redis.Do(re.context, "DEL", key).Result()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (re *redisStorage) SaveOrCreate(key string, age int, value []byte) error {
	_, err := re.redis.Do(re.context, "SETEX", key, age, value).Text()
	if err != nil {
		return err
	}
	return nil
}

func (re *redisStorage) Delete(key string) error {
	_, err := re.redis.Do(re.context, "DEL", key).Text()
	if err != nil {
		return err
	}
	return nil
}
