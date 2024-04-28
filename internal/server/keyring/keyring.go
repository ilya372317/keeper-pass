package keyring

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/cryptomanager"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

const generalKeySize = 32

var cryptoKey []byte

type keyStorage interface {
	GetKey(context.Context) (*domain.Keys, error)
	SaveKey(ctx context.Context, keys *domain.Keys) error
}

// Keyring is a keyring for general key and set of functions for manipulate with them.
type Keyring struct {
	keyStorage keyStorage
	masterKey  []byte
}

// New creates new Keyring instance.
func New(masterKey []byte, keyStorage keyStorage) *Keyring {
	return &Keyring{
		keyStorage: keyStorage,
		masterKey:  masterKey,
	}
}

// GetGeneralKey returns the general key in open form.
func (k *Keyring) GetGeneralKey(ctx context.Context) ([]byte, error) {
	if cryptoKey == nil {
		if err := k.InitGeneralKey(ctx); err != nil {
			return nil, fmt.Errorf("failed init general key: %w", err)
		}
	}

	return cryptoKey, nil
}

// InitGeneralKey initializes the general key.
func (k *Keyring) InitGeneralKey(ctx context.Context) error {
	aesgcm, err := cryptomanager.NewAESGCM(k.masterKey)
	if err != nil {
		return fmt.Errorf("failed init general key: %w", err)
	}
	keys, err := k.keyStorage.GetKey(ctx)
	if err == nil {
		var generalKey []byte
		generalKey, err = hex.DecodeString(keys.Key)
		if err != nil {
			return fmt.Errorf("failed decode key from storage: %w", err)
		}
		var nonce []byte
		nonce, err = hex.DecodeString(keys.Nonce)
		if err != nil {
			return fmt.Errorf("failed decode nonce from storage: %w", err)
		}

		var openKey []byte
		openKey, err = cryptomanager.Decrypt(aesgcm, generalKey, nonce)
		if err != nil {
			return fmt.Errorf("failed open general key: %w", err)
		}
		cryptoKey = openKey

		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed get key from storage: %w", err)
	}

	generalKey, err := cryptomanager.GenerateRandom(generalKeySize)
	if err != nil {
		return fmt.Errorf("failed generate random key: %w", err)
	}
	cryptoKey = generalKey

	nonce, err := cryptomanager.GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		return fmt.Errorf("failed generate nonce: %w", err)
	}
	keyForSave := cryptomanager.Encrypt(aesgcm, generalKey, nonce)

	key := &domain.Keys{
		Key:   hex.EncodeToString(keyForSave),
		Nonce: hex.EncodeToString(nonce),
	}

	if err = k.keyStorage.SaveKey(ctx, key); err != nil {
		return fmt.Errorf("failed save key to storage: %w", err)
	}

	return nil
}
