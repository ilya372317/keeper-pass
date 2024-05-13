package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
)

type AuthService interface {
	Login(context.Context, dto.LoginDTO) (string, error)
	Register(context.Context, dto.RegisterDTO) (*domain.User, error)
}

type loginPassService interface {
	Save(ctx context.Context, d dto.SaveLoginPassDTO) error
	Update(ctx context.Context, d dto.UpdateLoginPassDTO) error
	Show(ctx context.Context, id int) (domain.LoginPass, error)
}

type creditCardService interface {
	Save(ctx context.Context, d dto.SaveCreditCardDTO) error
	Update(ctx context.Context, d dto.UpdateCreditCardDTO) error
	Show(ctx context.Context, id int64) (domain.CreditCard, error)
}

type textService interface {
	Save(ctx context.Context, d dto.SaveTextDTO) error
	Update(ctx context.Context, d dto.UpdateTextDTO) error
	Show(ctx context.Context, id int64) (domain.Text, error)
}

type binaryService interface {
	Save(ctx context.Context, d dto.SaveBinaryDTO) error
	Update(ctx context.Context, d dto.UpdateBinaryDTO) error
	Show(ctx context.Context, id int64) (domain.Binary, error)
}

type generalDataService interface {
	Index(ctx context.Context) ([]domain.GeneralData, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	pb.UnimplementedPassServiceServer
	authService        AuthService
	loginPassService   loginPassService
	creditCardService  creditCardService
	textService        textService
	binaryService      binaryService
	generalDataService generalDataService
}

func New(
	authService AuthService,
	loginPassService loginPassService,
	creditCardService creditCardService,
	textService textService,
	binaryService binaryService,
	generalDataService generalDataService,
) *Handler {
	return &Handler{
		authService:        authService,
		loginPassService:   loginPassService,
		creditCardService:  creditCardService,
		textService:        textService,
		binaryService:      binaryService,
		generalDataService: generalDataService,
	}
}
