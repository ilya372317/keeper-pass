package cryptomanager

import (
	"crypto/cipher"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAESGCM(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "Valid key length 16 bytes",
			key:     make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "Valid key length 24 bytes",
			key:     make([]byte, 24),
			wantErr: false,
		},
		{
			name:    "Valid key length 32 bytes",
			key:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "Invalid key length 15 bytes",
			key:     make([]byte, 15),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aesgcm, err := NewAESGCM(tt.key)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, aesgcm)
			} else {
				require.NoError(t, err)
				require.NotNil(t, aesgcm)
			}
		})
	}
}

func TestGenerateRandom(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Generate 16 bytes",
			size:    16,
			wantErr: false,
		},
		{
			name:    "Generate 32 bytes",
			size:    32,
			wantErr: false,
		},
		{
			name:    "Generate 0 bytes",
			size:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandom(tt.size)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, got, tt.size, "the generated slice should have the requested size")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	key, err := GenerateRandom(32)
	require.NoError(t, err)

	aesgcm, err := NewAESGCM(key)
	require.NoError(t, err)

	plaintext := []byte("secret message")
	nonce, err := GenerateRandom(aesgcm.NonceSize())
	require.NoError(t, err)

	encryptedData := Encrypt(aesgcm, plaintext, nonce)

	tests := []struct {
		name    string
		aesgcm  cipher.AEAD
		data    []byte
		nonce   []byte
		wantErr bool
	}{
		{
			name:    "Successful decryption",
			aesgcm:  aesgcm,
			data:    encryptedData,
			nonce:   nonce,
			wantErr: false,
		},
		{
			name:    "Decryption with incorrect nonce",
			aesgcm:  aesgcm,
			data:    encryptedData,
			nonce:   []byte("wrongnonce"),
			wantErr: true,
		},
		{
			name:    "Decryption with tampered data",
			aesgcm:  aesgcm,
			data:    append(encryptedData, byte(1)),
			nonce:   nonce,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Errorf("TestDecrypt panicked for %s, but no panic expected. Panic: %v", tt.name, r)
					}
				}
			}()
			_, err := Decrypt(tt.aesgcm, tt.data, tt.nonce)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
