package database

import (
	"errors"
	"regexp"
)

func ValidateEmail(e string) error {
	if len(e) > 60 {
		return errors.New("Email address is too long")
	}
	matched, err := regexp.Match("^.+@.+\\..+$", []byte(e))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("Invalid email address")
	}
	return nil
}

func ValidatePassword(p string) error {
	if len(p) < 8 || len(p) > 60 {
		return errors.New("Password length is invalid (8-60)")
	}
	return nil
}
