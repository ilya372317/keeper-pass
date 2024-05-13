package binary

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	binary_mock "github.com/ilya372317/pass-keeper/internal/server/service/binary/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := binary_mock.NewMockdataService(ctrl)
	serv := Service{dataService: dataServ}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{ID: 1})

	t.Run("success update case", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateBinaryDTO{
			Metadata: &dto.BinaryMetadata{Info: "info"},
			Data:     byteSlicePtr([]byte("132")),
			ID:       1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"MTMy"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindBinary,
			IsPayloadDecrypted: true,
		}, nil)
		dataServ.EXPECT().EncryptAndUpdateData(gomock.Any(), domain.Data{
			Payload:            `{"data":"MTMy"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindBinary,
			IsPayloadDecrypted: true,
		})

		// Execute.
		err := serv.Update(ctxUser, arg)

		// Assert.
		require.NoError(t, err)
	})

	t.Run("missing user in ctx", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateBinaryDTO{ID: 1}

		// Execute.
		err := serv.Update(context.Background(), arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed get or decrypt binary data", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateBinaryDTO{ID: 1}
		dataServ.EXPECT().
			GetAndDecryptData(gomock.Any(), int64(1)).
			Times(1).
			Return(domain.Data{}, fmt.Errorf("internal"))

		// Execute.
		err := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("attempt to update binary data belongs to another user", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateBinaryDTO{
			ID: 1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(
			domain.Data{
				ID:                 1,
				UserID:             2,
				Kind:               domain.KindBinary,
				IsPayloadDecrypted: true,
			}, nil)

		// Execute.
		err := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAccesDenied)
	})

	t.Run("invalid kind", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateBinaryDTO{
			ID: 1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		err := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotSupportedOperation)
	})

	invalidDataInStorageTests := []struct {
		name     string
		payload  string
		metadata string
	}{
		{
			name:     "metadata",
			payload:  `{"data":"MTMy"}`,
			metadata: "invalid metadata",
		},
		{
			name:     "payload",
			payload:  `invalid-payload`,
			metadata: `{"info":"info"}`,
		},
	}

	for _, tt := range invalidDataInStorageTests {
		t.Run("invalid "+tt.name+" in storage", func(t *testing.T) {
			// Prepare.
			arg := dto.UpdateBinaryDTO{
				ID: 1,
			}
			dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
				Payload:            tt.payload,
				Metadata:           tt.metadata,
				ID:                 1,
				UserID:             1,
				Kind:               domain.KindBinary,
				IsPayloadDecrypted: true,
			}, nil)

			// Execute.
			err := serv.Update(ctxUser, arg)

			// Assert.
			require.Error(t, err)
		})
	}
}

func byteSlicePtr(val []byte) *[]byte {
	return &val
}
