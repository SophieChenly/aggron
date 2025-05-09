package services

import (
	"aggron/internal/db/models"
	"aggron/internal/repository"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileKey struct {
	FileID            string
	EncryptedKey      []byte // KMS-encrypted key
	FileHash          []byte // SHA-256 hash of the original file
	SenderDiscordID   string
	ReceiverDiscordID string
}

type FileEncryptionService struct {
	encryptionService *EncryptionService
	kmsService        *KMSKeyService
	fileKeyRepo       *repository.FileKeyRepository
}

func NewFileEncryptionService(
	encryptionService *EncryptionService,
	kmsService *KMSKeyService,
	fileKeyRepo *repository.FileKeyRepository,
) *FileEncryptionService {
	return &FileEncryptionService{
		encryptionService: encryptionService,
		kmsService:        kmsService,
		fileKeyRepo:       fileKeyRepo,
	}
}

// EncryptFile encrypts a file and stores the encrypted key
func (s *FileEncryptionService) EncryptFile(
	ctx context.Context,
	fileID string,
	fileData []byte,
	senderDiscordID string,
	receiverDiscordID string,
) ([]byte, error) {
	// checksum
	fileHash := s.encryptionService.GenerateFileHash(fileData)

	plaintextKey, encryptedKey, err := s.kmsService.GenerateDataKey(ctx)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(plaintextKey)

	encryptedData, err := s.encryptionService.Encrypt(fileData, plaintextKey)
	if err != nil {
		return nil, err
	}

	fileKey := models.FileKey{
		FileID:            fileID,
		EncryptedKey:      primitive.Binary{Data: encryptedKey},
		FileHash:          primitive.Binary{Data: fileHash},
		SenderDiscordID:   senderDiscordID,
		ReceiverDiscordID: receiverDiscordID,
		CreatedAt:         time.Now(),
	}

	_, err = s.fileKeyRepo.CreateFileKey(ctx, fileKey)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func (s *FileEncryptionService) DecryptFile(
	ctx context.Context,
	fileID string,
	encryptedData []byte,
	userDiscordID string,
) ([]byte, error) {

	fileKey, err := s.fileKeyRepo.FindByFileID(ctx, fileID)
	if err != nil {
		return nil, errors.New("file key not found")
	}

	if fileKey.SenderDiscordID != userDiscordID && fileKey.ReceiverDiscordID != userDiscordID {
		return nil, errors.New("user does not have permission to access this file")
	}

	plaintextKey, err := s.kmsService.DecryptDataKey(ctx, fileKey.EncryptedKey.Data)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(plaintextKey)

	decryptedData, err := s.encryptionService.Decrypt(encryptedData, plaintextKey)
	if err != nil {
		return nil, err
	}

	if !s.encryptionService.ValidateFileHash(decryptedData, fileKey.FileHash.Data) {
		return nil, errors.New("file integrity verification failed")
	}

	return decryptedData, nil
}

// overwrites a byte slice with zeros in memory
func zeroBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
