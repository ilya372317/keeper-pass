package dto

type BinaryMetadata struct {
	Info string `json:"info,omitempty"`
}

type SaveBinaryDTO struct {
	Metadata TextMetadata `json:"-"`
	Data     []byte       `json:"data"`
}

type UpdateBinaryDTO struct {
	Metadata *TextMetadata
	Data     []byte
	ID       int64
}

type BinaryPayload struct {
	Data []byte `json:"data"`
}
