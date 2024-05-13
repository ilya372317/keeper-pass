package binary

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	binary_mock "github.com/ilya372317/pass-keeper/internal/server/service/binary/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := binary_mock.NewMockdataService(ctrl)
	serv := Service{dataService: dataServ}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{ID: 1})

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"MTMy"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindBinary,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		got, err := serv.Show(ctxUser, 1)

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, []byte("132"), got.Data)
		assert.Equal(t, "info", got.Metadata.Info)
		assert.Equal(t, int64(1), got.ID)
	})

	t.Run("missing user in ctx", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()

		// Execute.
		_, err := serv.Show(ctx, 1)

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed get data from data service", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{},
			fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
	})

	errTests := []struct {
		name string
		data domain.Data
		want error
	}{
		{
			name: "permission denied",
			data: domain.Data{
				ID:                 1,
				UserID:             2,
				Kind:               domain.KindBinary,
				IsPayloadDecrypted: false,
			},
			want: domain.ErrAccesDenied,
		},
		{
			name: "invalid kind",
			data: domain.Data{
				ID:                 1,
				UserID:             1,
				Kind:               domain.KindText,
				IsPayloadDecrypted: true,
			},
			want: domain.ErrNotSupportedOperation,
		},
	}
	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(tt.data, nil)

			// Execute.
			_, err := serv.Show(ctxUser, 1)

			// Assert.
			require.Error(t, err)
			require.ErrorIs(t, err, tt.want)
		})
	}

}
