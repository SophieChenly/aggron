package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	SetObjMap(ctx context.Context, key string, value any, expiration time.Duration) error
	GetObjMap(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type Redis struct {
	client redis.Client
}

func NewRedis(client redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func SetObjTyped[T any](r *Redis, ctx context.Context, key string, value T, expiration time.Duration) error {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, b.Bytes(), expiration).Err()
}

func GetObjTyped[T any](r *Redis, ctx context.Context, key string) (T, error) {
	var empty T

	result := r.client.Get(ctx, key)
	cmdb, err := result.Bytes()
	if err != nil {
		return empty, err
	}

	b := bytes.NewReader(cmdb)

	var res T
	err = gob.NewDecoder(b).Decode(&res)
	if err != nil {
		return empty, err
	}

	return res, nil
}

func (r *Redis) GetObjMap(ctx context.Context, key string) (map[string]string, error) {
	result := r.client.Get(ctx, key)
	cmdb, err := result.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res map[string]string
	err = gob.NewDecoder(b).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Redis) SetObjMap(ctx context.Context, key string, value map[string]string, expiration time.Duration) error {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, b.Bytes(), expiration).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
