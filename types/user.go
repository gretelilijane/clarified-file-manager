package types

type User struct {
	ID           uint64
	Username     string
	PasswordHash []byte
	PasswordSalt []byte
}
