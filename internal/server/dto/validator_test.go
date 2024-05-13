package dto

import (
	"testing"
)

func TestDTO_ValidateWithRequired(t *testing.T) {
	type fields struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		wantErr bool
		fields  fields
	}{
		{
			name:    "success case",
			wantErr: false,
			fields: fields{
				email:    "ilya.otinov@gmail.com",
				password: "123",
			},
		},
		{
			name:    "incorrect email case",
			wantErr: true,
			fields: fields{
				email:    "email-invalid",
				password: "123",
			},
		},
		{
			name:    "incorrect password case",
			wantErr: true,
			fields: fields{
				email:    "ilya.otinov@gmail.com",
				password: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &LoginDTO{
				Email:    tt.fields.email,
				Password: tt.fields.password,
			}
			if err := ValidateDTOWithRequired(d); (err != nil) != tt.wantErr {
				t.Errorf("ValidateWithRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
