package dto

type LoginPassMetadata struct {
	URL string `json:"url,omitempty" validate:"omitempty,url"`
}
