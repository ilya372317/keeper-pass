package data

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/cryptomanager"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) SaveSimpleData(ctx context.Context, d dto.SaveSimpleDataDTO) (*domain.Data, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return nil, fmt.Errorf("failed get user from context")
	}

	if err := checkPayload(d); err != nil {
		return nil, fmt.Errorf("failed validate payload: %w", err)
	}

	cryptokey, err := cryptomanager.GenerateRandom(domain.CryptoKeyLength)
	if err != nil {
		return nil, fmt.Errorf("failed generate random key: %w", err)
	}
	aesgcmPayload, err := cryptomanager.NewAESGCM(cryptokey)
	if err != nil {
		return nil, fmt.Errorf("failed create new aes gcm: %w", err)
	}
	noncePayload, err := cryptomanager.GenerateRandom(aesgcmPayload.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("failed generate noncePayload: %w", err)
	}

	cryptedPayload := cryptomanager.Encrypt(aesgcmPayload, []byte(d.Payload), noncePayload)

	aesgcmDataKey, err := cryptomanager.NewAESGCM(cryptokey)
	if err != nil {
		return nil, fmt.Errorf("failed create new aes gcm: %w", err)
	}

	nonceDataKey, err := cryptomanager.GenerateRandom(aesgcmDataKey.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("failed generate nonceDataKey: %w", err)
	}

	cryptedDataKey := cryptomanager.Encrypt(aesgcmDataKey, cryptokey, nonceDataKey)

	data := domain.Data{
		Payload:   hex.EncodeToString(cryptedPayload),
		Metadata:  d.Metadata,
		Nonce:     hex.EncodeToString(noncePayload),
		CryptoKey: hex.EncodeToString(cryptedDataKey),
		UserID:    int(user.ID),
		Kind:      d.Type,
	}

	return &data, nil
}

func checkPayload(d dto.SaveSimpleDataDTO) error {
	switch d.Type {
	case domain.KindFile:
		return errors.Join(fmt.Errorf("file saving not supported by unary operation"), domain.ErrPayloadNotValid)
	case domain.KindLoginPass:
		payload := dto.LoginPassPayload{}
		if err := json.Unmarshal([]byte(d.Payload), &payload); err != nil {
			return errors.Join(err, domain.ErrPayloadNotValid)
		}
		if err := dto.ValidateDTOWithRequired(&payload); err != nil {
			return errors.Join(err, domain.ErrPayloadNotValid)
		}
	case domain.KindCreditCard:
		payload := dto.CreditCardPayload{}
		if err := json.Unmarshal([]byte(d.Payload), &payload); err != nil {
			return errors.Join(err, domain.ErrPayloadNotValid)
		}
		if err := dto.ValidateDTOWithRequired(&payload); err != nil {
			return errors.Join(err, domain.ErrPayloadNotValid)
		}
	}

	return nil
}
