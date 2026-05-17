package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const cost = bcrypt.DefaultCost

func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", fmt.Errorf("password: hash: %w", err)
	}
	return string(b), nil
}

func Compare(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
