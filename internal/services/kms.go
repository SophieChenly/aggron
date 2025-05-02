package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KMSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Timeout         time.Duration
}

type KMSKeyService struct {
	client *kms.Client
	keyID  string
}

func NewKMSKeyService(config KMSConfig, keyID string) (*KMSKeyService, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	creds := credentials.NewStaticCredentialsProvider(
		config.AccessKeyID,
		config.SecretAccessKey,
		"", // optional session token
	)

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(config.Region),
		awsconfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := kms.NewFromConfig(cfg)
	return &KMSKeyService{
		client: client,
		keyID:  keyID,
	}, nil
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
