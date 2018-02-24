package database

import (
	"errors"
	"regexp"
)

var ValidEmailRegExp = "^(?=[^]{,60}$)[^]+@[^]+\\.[^]+$"
var ValidPasswordRegExp = "^[^]{8,60}$"

func validateEmail(e string) error {
	matched, err := regexp.Match(ValidEmailRegExp, []byte(e))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("Invalid email address")
	}
	return nil
}

func validatePassword(p string) error {
	matched, err := regexp.Match(ValidPasswordRegExp, []byte(p))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("Password length is invalid (8-60)")
	}
	return nil
}
