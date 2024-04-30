package domain

// Kind represent a data type.
type Kind int

const (
	KindLoginPass  Kind = iota // KindLoginPass represent a login with password data type.
	KindFile                   // KindFile represent a file data type.
	KindCreditCard             // KindCreditCard represent a credit card data type.
)

// Data represent a storing data.
type Data struct {
	Payload   string `db:"payload"`    // Payload of the data.
	Metadata  string `db:"metadata"`   // Metadata of the data.
	Nonce     string `db:"nonce"`      // Nonce of the data.
	CryptoKey string `db:"crypto_key"` // Crypto key of the data.
	ID        int    `db:"id"`         // ID of the data.
	Kind      Kind   `db:"kind"`       // Kind of the data.
}
