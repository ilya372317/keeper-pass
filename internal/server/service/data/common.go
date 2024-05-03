package data

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/cryptomanager"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) encryptData(ctx context.Context, data domain.Data) (domain.Data, error) {
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

	cryptedPayload := cryptomanager.Encrypt(aesgcmPayload, []byte(data.Payload), noncePayload)

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
		Payload:            hex.EncodeToString(cryptedPayload),
		Metadata:           data.Metadata,
		PayloadNonce:       hex.EncodeToString(noncePayload),
		CryptoKey:          hex.EncodeToString(cryptedDataKey),
		CryptoKeyNonce:     hex.EncodeToString(nonceDataKey),
		UserID:             data.UserID,
		Kind:               data.Kind,
		IsPayloadDecrypted: false,
	}, nil
}

func (s *Service) decryptData(ctx context.Context, data domain.Data) (domain.Data, error) {
	generalKey, err := s.keyring.GetGeneralKey(ctx)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed get general key from storage: %w", err)
	}

	dataKeyAesgcm, err := cryptomanager.NewAESGCM(generalKey)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed create new aesgcm for decrypt data key: %w", err)
	}

	dataKeyNonce, err := hex.DecodeString(data.CryptoKeyNonce)
	if err != nil {
		return domain.Data{}, fmt.Errorf("crypto key nonce expected to be in hex codec: %w", err)
	}

	cryptedDataKey, err := hex.DecodeString(data.CryptoKey)
	if err != nil {
		return domain.Data{}, fmt.Errorf("crypto key expected to be in hex codec: %w", err)
	}

	decryptedDataKey, err := dataKeyAesgcm.Open(nil, dataKeyNonce, cryptedDataKey, nil)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed decrypt data key: %w", err)
	}

	payloadAesgcm, err := cryptomanager.NewAESGCM(decryptedDataKey)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed create aesgcm for decrypt data payload")
	}

	payloadNonce, err := hex.DecodeString(data.PayloadNonce)
	if err != nil {
		return domain.Data{}, fmt.Errorf("payload nonce expected to be in hex codec: %w", err)
	}

	cryptedPayload, err := hex.DecodeString(data.Payload)
	if err != nil {
		return domain.Data{}, fmt.Errorf("payload expected to be in hex codec: %w", err)
	}

	decryptedPayload, err := payloadAesgcm.Open(nil, payloadNonce, cryptedPayload, nil)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed decrypt data payload: %w", err)
	}

	return domain.Data{
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
		Payload:            string(decryptedPayload),
		Metadata:           data.Metadata,
		PayloadNonce:       data.PayloadNonce,
		CryptoKeyNonce:     data.CryptoKeyNonce,
		CryptoKey:          data.CryptoKey,
		ID:                 data.ID,
		UserID:             data.UserID,
		Kind:               data.Kind,
		IsPayloadDecrypted: true,
	}, nil
}
