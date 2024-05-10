package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinary_GetInfo(t *testing.T) {
	// Prepare.
	binary := Binary{Metadata: BinaryMetadata{Info: "info"}}

	// Execute.
	got := binary.GetInfo()

	// Assert.
	assert.Equal(t, "info", got)
}

func TestCreditCard_GetInfo(t *testing.T) {
	// Prepare.
	cc := CreditCard{Metadata: CreditCardMetadata{BankName: "info"}}

	// Execute.
	got := cc.GetInfo()

	// Assert.
	assert.Equal(t, "info", got)
}

func TestText_GetInfo(t *testing.T) {
	// Prepare.
	text := Text{Metadata: TextMetadata{Info: "info"}}

	// Execute.
	got := text.GetInfo()

	// Assert.
	assert.Equal(t, "info", got)
}

func TestLoginPass_GetInfo(t *testing.T) {
	// Prepare.
	lp := LoginPass{Metadata: LoginPassMetadata{URL: "info"}}

	// Execute.
	got := lp.GetInfo()

	// Assert.
	assert.Equal(t, "info", got)
}

func TestData_ToDataRepresentation(t *testing.T) {
	type want struct {
		res string
		err bool
	}
	type arg struct {
		metadata string
		kind     Kind
	}
	tests := []struct {
		name string
		want want
		arg  arg
	}{
		{
			name: "login pass case",
			want: want{
				res: "info",
				err: false,
			},
			arg: arg{
				metadata: `{"URL":"info"}`,
				kind:     KindLoginPass,
			},
		},
		{
			name: "credit card case",
			want: want{
				res: "info",
				err: false,
			},
			arg: arg{
				metadata: `{"bank_name":"info"}`,
				kind:     KindCreditCard,
			},
		},
		{
			name: "text case",
			want: want{
				res: "info",
				err: false,
			},
			arg: arg{
				metadata: `{"info":"info"}`,
				kind:     KindText,
			},
		},
		{
			name: "binary case",
			want: want{
				res: "info",
				err: false,
			},
			arg: arg{
				metadata: `{"info":"info"}`,
				kind:     KindBinary,
			},
		},
		{
			name: "unknown kind",
			want: want{
				res: "",
				err: true,
			},
			arg: arg{
				metadata: `{}`,
				kind:     -1,
			},
		},
		{
			name: "invalid login pass metadata",
			want: want{
				res: "",
				err: true,
			},
			arg: arg{
				metadata: `invalid login pass metadata`,
				kind:     KindLoginPass,
			},
		},
		{
			name: "invalid credit card metadata",
			want: want{
				res: "",
				err: true,
			},
			arg: arg{
				metadata: `"not_bank_name":"info"`,
				kind:     KindCreditCard,
			},
		},
		{
			name: "invalid text metadata",
			want: want{
				res: "",
				err: true,
			},
			arg: arg{
				metadata: `"info":"info"`,
				kind:     KindText,
			},
		},
		{
			name: "invalid binary metadata",
			want: want{
				res: "info",
				err: true,
			},
			arg: arg{
				metadata: `"info":"info"`,
				kind:     KindBinary,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			data := Data{Metadata: tt.arg.metadata, Kind: tt.arg.kind}

			// Execute.
			got, err := data.ToDataRepresentation()

			// Assert.
			if tt.want.err {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, got)

			assert.Equal(t, tt.want.res, got.GetInfo())
		})
	}
}
