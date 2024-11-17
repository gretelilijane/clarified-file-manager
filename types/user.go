package types

type User struct {
	ID           int32
	Username     string
	PasswordHash []byte
	PasswordSalt []byte
}
