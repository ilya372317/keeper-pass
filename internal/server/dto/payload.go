package dto

type LoginPassPayloadDTO struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
