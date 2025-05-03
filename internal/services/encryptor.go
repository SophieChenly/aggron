package services

import (
	"context"
	"errors"
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
	keyStore          *KeyStoreService
}

func NewFileEncryptionService(
	encryptionService *EncryptionService,
	kmsService *KMSKeyService,
	keyStore *KeyStoreService,
) *FileEncryptionService {
	return &FileEncryptionService{
		encryptionService: encryptionService,
		kmsService:        kmsService,
		keyStore:          keyStore,
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

	fileKey := FileKey{
		FileID:            fileID,
		EncryptedKey:      encryptedKey,
		FileHash:          fileHash,
		SenderDiscordID:   senderDiscordID,
		ReceiverDiscordID: receiverDiscordID,
	}

	if err := s.keyStore.StoreFileKey(fileKey); err != nil {
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

	fileKey, err := s.keyStore.GetFileKey(fileID)
	if err != nil {
		return nil, errors.New("file key not found")
	}

	if fileKey.SenderDiscordID != userDiscordID && fileKey.ReceiverDiscordID != userDiscordID {
		return nil, errors.New("user does not have permission to access this file")
	}

	plaintextKey, err := s.kmsService.DecryptDataKey(ctx, fileKey.EncryptedKey)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(plaintextKey)

	decryptedData, err := s.encryptionService.Decrypt(encryptedData, plaintextKey)
	if err != nil {
		return nil, err
	}

	if !s.encryptionService.ValidateFileHash(decryptedData, fileKey.FileHash) {
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
