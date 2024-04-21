package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestUser_IsPasswordCorrect(t *testing.T) {
	type fields struct {
		password string
		salt     string
	}
	tests := []struct {
		name   string
		fields fields
		arg    string
		want   bool
	}{
		{
			name: "success correct passwords case",
			fields: fields{
				password: "123",
				salt:     "my-salt",
			},
			arg:  "123",
			want: true,
		},
		{
			name: "invalid password case",
			fields: fields{
				password: "123",
				salt:     "salt",
			},
			arg:  "321",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Salt: tt.fields.salt,
			}
			u.SetHashedPassword(tt.fields.password)

			got := u.IsPasswordCorrect(tt.arg)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUser_SetHashedPassword(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		salt string
	}{
		{
			name: "simple not long case",
			arg:  "123",
			salt: "salt-is-good",
		},
		{
			name: "long salt case",
			arg:  "123",
			salt: "salt-is-good and in this string we have spaces and it longer then 32 bytes." +
				" 32 bytes is a length of result excecution hash function. And it is very funny." +
				" Maybe i don`t need to spent this time to write it, but i have fun and it is only thing metter.",
		},
		{
			name: "long password case",
			arg: "password-is-good and in this string we have spaces and it longer then 32 bytes." +
				" 32 bytes is a length of result excecution hash function. And it is very funny." +
				" Maybe i don`t need to spent this time to write it, but i have fun and it is only thing metter.",
			salt: "123",
		},
		{
			name: "empty password case",
			arg:  "",
			salt: "salt",
		},
		{
			name: "empty salt case",
			arg:  "123-pass",
			salt: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Salt: tt.salt,
			}

			u.SetHashedPassword(tt.arg)

			got := u.HashedPassword

			argWithSaltBytes := sha256.Sum256([]byte(tt.arg + tt.salt))
			argWithSalt := hex.EncodeToString(argWithSaltBytes[:])

			assert.Equal(t, argWithSalt, got)
		})
	}
}
