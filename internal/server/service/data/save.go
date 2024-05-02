package data

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/cryptomanager"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) EncryptAndSaveData(ctx context.Context, d dto.SaveSimpleDataDTO) error {
	data, err := s.buildCryptedData(ctx, d)
	if err != nil {
		return fmt.Errorf("failed encrypt data: %w", err)
	}

	if err = s.dataStorage.SaveData(ctx, data); err != nil {
		return fmt.Errorf("failed save data to storage: %w", err)
	}

	return nil
}

func (s *Service) buildCryptedData(ctx context.Context, d dto.SaveSimpleDataDTO) (domain.Data, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return domain.Data{}, fmt.Errorf("failed get user from context")
	}

	cryptokey, err := cryptomanager.GenerateRandom(domain.CryptoKeyLength)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed generate random key: %w", err)
	}
	aesgcmPayload, err := cryptomanager.NewAESGCM(cryptokey)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed create new aes gcm: %w", err)
	}
	noncePayload, err := cryptomanager.GenerateRandom(aesgcmPayload.NonceSize())
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed generate noncePayload: %w", err)
	}

	cryptedPayload := cryptomanager.Encrypt(aesgcmPayload, []byte(d.Payload), noncePayload)

	generalKey, err := s.keyring.GetGeneralKey(ctx)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed get general key from keyring: %w", err)
	}
	aesgcmDataKey, err := cryptomanager.NewAESGCM(generalKey)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed create new aes gcm: %w", err)
	}

	nonceDataKey, err := cryptomanager.GenerateRandom(aesgcmDataKey.NonceSize())
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed generate nonceDataKey: %w", err)
	}

	cryptedDataKey := cryptomanager.Encrypt(aesgcmDataKey, cryptokey, nonceDataKey)

	return domain.Data{
		Payload:        hex.EncodeToString(cryptedPayload),
		Metadata:       d.Metadata,
		PayloadNonce:   hex.EncodeToString(noncePayload),
		CryptoKey:      hex.EncodeToString(cryptedDataKey),
		CryptoKeyNonce: hex.EncodeToString(nonceDataKey),
		UserID:         user.ID,
		Kind:           d.Type,
	}, nil
}
