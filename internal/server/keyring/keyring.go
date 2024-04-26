package keyring

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

const generalKeySize = 32

type keyStorage interface {
	GetKey(context.Context) (*domain.Keys, error)
	SaveKey(ctx context.Context, key string) error
}

type Keyring struct {
	keyStorage keyStorage
	GeneralKey []byte
}

func New(ctx context.Context, masterKey string, keyStorage keyStorage) (*Keyring, error) {
	keyRing := &Keyring{
		keyStorage: keyStorage,
	}

	keys, err := keyRing.keyStorage.GetKey(ctx)
	if err == nil {
		generalKey, err := hex.DecodeString(keys.Key)
		if err != nil {
			return nil, fmt.Errorf("failed decode key from storage: %w", err)
		}

		keyRing.GeneralKey = generalKey
		return keyRing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed get key from storage: %w", err)
	}

	generalKey, err := generateRandom(generalKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed generate random key: %w", err)
	}
	keyRing.GeneralKey = generalKey

	aesblock, err := aes.NewCipher([]byte(masterKey))
	if err != nil {
		return nil, fmt.Errorf("failed create new aes block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("failed create new aes gcm: %w", err)
	}

	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("failed generate random nonce: %w", err)
	}

	keyForSave := aesgcm.Seal(nil, nonce, generalKey, nil)

	if err = keyRing.keyStorage.SaveKey(ctx, hex.EncodeToString(keyForSave)); err != nil {
		return nil, fmt.Errorf("failed save key to storage: %w", err)
	}

	return keyRing, nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed generate random bytes: %w", err)
	}

	return b, nil
}
