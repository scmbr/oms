package hasher

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(pwd string) (string, error)
	Compare(hash, pwd string) bool
}
type BcryptHasher struct{}

func (BcryptHasher) Hash(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(b), err
}

func (BcryptHasher) Compare(hash, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
