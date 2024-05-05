package dto

type TextMetadata struct {
	Info string `json:"info,omitempty"`
}

type SaveTextDTO struct {
	Metadata TextMetadata `json:"-"`
	Data     string       `json:"data"`
}

type UpdateTextDTO struct {
	Metadata *TextMetadata
	Data     *string
	ID       int64
}

type TextPayload struct {
	Data string `json:"data"`
}
