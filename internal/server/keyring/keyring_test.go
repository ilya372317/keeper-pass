package keyring

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	keyringmock "github.com/ilya372317/pass-keeper/internal/server/keyring/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type getKeySettings struct {
		returnKey *domain.Keys
		returnErr error
	}
	type saveKeySettings struct {
		returnErr error
	}
	type keyStorageSettings struct {
		getKeySettings  getKeySettings
		saveKeySettings saveKeySettings
	}
	type want struct {
		err bool
	}
	tests := []struct {
		name               string
		masterKey          string
		keyStorageSettings keyStorageSettings
		want               want
	}{
		{
			name:      "storage has key case",
			masterKey: "key",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: &domain.Keys{
						ID:        1,
						Key:       hex.EncodeToString([]byte("123")),
						IsCurrent: true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					returnErr: nil,
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: false,
			},
		},
		{
			name:      "failed decode string case",
			masterKey: "123",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: &domain.Keys{
						ID:        1,
						Key:       "123",
						IsCurrent: true,
					},
					returnErr: nil,
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name:      "failed get key from storage case",
			masterKey: "123",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: nil,
					returnErr: fmt.Errorf("failed get key from storage"),
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name:      "empty master key case",
			masterKey: "",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: nil,
					returnErr: sql.ErrNoRows,
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name:      "success case with 16 length key size",
			masterKey: "1234567891234567",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: nil,
					returnErr: sql.ErrNoRows,
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: false,
			},
		},
		{
			name:      "success case with 24 length key size",
			masterKey: "123456789123456789123456",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: nil,
					returnErr: sql.ErrNoRows,
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: false,
			},
		},
		{
			name:      "success case with 32 length key size",
			masterKey: "11111111111111111111111111111111",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: nil,
					returnErr: sql.ErrNoRows,
				},
				saveKeySettings: saveKeySettings{
					returnErr: nil,
				},
			},
			want: want{
				err: false,
			},
		},
		{
			name:      "failed save generated key to storage case",
			masterKey: "1111111111111111",
			keyStorageSettings: keyStorageSettings{
				getKeySettings: getKeySettings{
					returnKey: nil,
					returnErr: sql.ErrNoRows,
				},
				saveKeySettings: saveKeySettings{
					returnErr: fmt.Errorf("failed save generated key to storage"),
				},
			},
			want: want{
				err: true,
			},
		},
	}
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := keyringmock.NewMockkeyStorage(ctrl)
			storage.
				EXPECT().
				SaveKey(gomock.Any(), gomock.Any()).
				Return(tt.keyStorageSettings.saveKeySettings.returnErr).
				AnyTimes()
			storage.
				EXPECT().
				GetKey(gomock.Any()).
				AnyTimes().
				Return(tt.keyStorageSettings.getKeySettings.returnKey, tt.keyStorageSettings.getKeySettings.returnErr)

			got, err := New(ctx, tt.masterKey, storage)
			if tt.want.err {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.True(t, len(got.GeneralKey) > 0)
		})
	}
}
