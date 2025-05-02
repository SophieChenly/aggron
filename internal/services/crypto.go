package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

type EncryptionService struct {
	keySize int
}

func NewEncryptionService() *EncryptionService {
	return &EncryptionService{
		keySize: 32, // AES-256
	}
}

func (s *EncryptionService) GenerateKey() ([]byte, error) {
	key := make([]byte, s.keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func (s *EncryptionService) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	// Combination of nonce + key MUST be unique, else catastrophic failure
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Ciphertext format: nonce + ciphertext
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, err
}

func (s *EncryptionService) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != s.keySize {
		return nil, errors.New("invalid key size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("malformed ciphertext")
	}

	var nonce, text = ciphertext[:nonceSize], ciphertext[nonceSize:]

	return aesgcm.Open(nil, nonce, text, nil)
}

func (s *EncryptionService) GenerateFileHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (s *EncryptionService) ValidateFileHash(data []byte, expectedHash []byte) bool {
	hash := sha256.Sum256(data)
	return string(hash[:]) == string(expectedHash)
}
