package security

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type bcryptHasher struct {
	cost int
}

var _ PasswordHasher = (*bcryptHasher)(nil)

func NewBcryptHasher(cost int) *bcryptHasher {
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

// Returns nil on success, or an error on failure
func (h *bcryptHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
