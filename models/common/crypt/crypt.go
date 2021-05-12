package crypt

import (
	"Cube-back/log"
	"golang.org/x/crypto/bcrypt"
)

func Set(word string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(word), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
	}
	return string(hash)
}

func Confirm(compare, origin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(origin), []byte(compare))
	if err != nil {
		return false
	} else {
		return true
	}
}
