package dto

import "github.com/ilya372317/pass-keeper/internal/server/domain"

type SaveSimpleDataDTO struct {
	Payload  string      `validate:"json"`
	Metadata string      `validate:"json"`
	Type     domain.Kind `validate:"min=0,max=2"`
}
