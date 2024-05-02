package dto

type SaveLoginPassDTO struct {
	Metadata LoginPassMetadata `json:"metadata,omitempty"`
	Login    string            `json:"login" validate:"min=3,max=255"`
	Password string            `json:"password" validate:"min=3,max=255"`
}
