package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KMSKeyService struct {
	client *kms.Client
	keyID  string
}

func NewKMSKeyService(cfg aws.Config, keyID string) *KMSKeyService {
	client := kms.NewFromConfig(cfg)
	return &KMSKeyService{
		client: client,
		keyID:  keyID,
	}
}

func (s *KMSKeyService) GenerateDataKey(ctx context.Context) ([]byte, []byte, error) {
	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(s.keyID),
		KeySpec: types.DataKeySpecAes256,
	}

	result, err := s.client.GenerateDataKey(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return result.Plaintext, result.CiphertextBlob, nil
}

func (s *KMSKeyService) DecryptDataKey(ctx context.Context, encryptedKey []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		CiphertextBlob: encryptedKey,
	}

	result, err := s.client.Decrypt(ctx, input)
	if err != nil {
		return nil, err
	}

	return result.Plaintext, nil
}

func (s *KMSKeyService) ReEncryptDataKey(ctx context.Context, encryptedKey []byte, destinationKeyID string) ([]byte, error) {
	input := &kms.ReEncryptInput{
		CiphertextBlob:   encryptedKey,
		SourceKeyId:      aws.String(s.keyID),
		DestinationKeyId: aws.String(destinationKeyID),
	}

	result, err := s.client.ReEncrypt(ctx, input)
	if err != nil {
		return nil, err
	}

	return result.CiphertextBlob, nil
}
