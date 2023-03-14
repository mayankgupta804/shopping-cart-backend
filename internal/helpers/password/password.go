package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const HASHING_COST = 10

var (
	ErrPasswordTooLong = errors.New("password is too long")
	ErrUnexpected      = errors.New("unexpected error encountered")
)

func GenerateHashedPassword(password []byte) ([]byte, error) {
	password, err := bcrypt.GenerateFromPassword(password, HASHING_COST)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return nil, ErrPasswordTooLong
	} else if err != nil {
		return nil, ErrUnexpected
	}
	return password, nil
}

func CompareHashAndPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		// maybe log how many bad logins happen in a given timeframe.
		// can it help in knowing if the system is being attacked?
		return false
	} else if err != nil {
		// log and report the error
		return false
	}
	return true
}
