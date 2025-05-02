package services

import (
	"aggron/internal/cache"
	"context"
	"errors"
)

type KeyStoreInterface interface {
	StoreFileKey(ctx context.Context, fileKey FileKey) error
	GetFileKey(ctx context.Context, fileID string) (FileKey, error)
	// TODO: Additional methods like ListFileKeys, DeleteFileKey, etc.
}

type KeyStoreService struct {
	redisService *cache.Redis
}

func NewKeyStoreService(redis *cache.Redis) *KeyStoreService {
	return &KeyStoreService{
		redisService: redis,
	}
}

func (k *KeyStoreService) StoreFileKey(fileKey FileKey) error {
	err := cache.SetObjTyped(k.redisService, context.TODO(), fileKey.FileID, fileKey, cache.DefaultExpirationTime)
	if err != nil {
		return errors.New("error setting filekey")
	}
	return nil
}

func (k *KeyStoreService) GetFileKey(fileID string) (FileKey, error) {
	fileKey, err := cache.GetObjTyped[FileKey](k.redisService, context.TODO(), fileID)
	if err != nil {
		return FileKey{}, errors.New("error getting filekey")
	}
	return fileKey, nil
}
