package dto

type SaveLoginPassDTO struct {
	Metadata LoginPassMetadata `json:"metadata,omitempty"`
	Login    string            `json:"login" validate:"min=3,max=255"`
	Password string            `json:"password" validate:"min=3,max=255"`
}

type UpdateLoginPassDTO struct {
	Metadata *LoginPassMetadata `json:"metadata" validate:"omitnil"`
	Login    *string            `json:"login" validate:"omitnil,min=3,max=255"`
	Password *string            `json:"password" validate:"omitnil,min=3,max=255"`
	ID       int64              `json:"id" validate:"required"`
}

type LoginPassMetadata struct {
	URL string `json:"url,omitempty" validate:"omitempty,url"`
}

type LoginPassPayloadDTO struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
