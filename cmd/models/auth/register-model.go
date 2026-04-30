package models_auth

import (
	"regexp"
	"unicode"
)

type Register struct {
	Username string
	Email    string
	Password string
}

func (r *Register) ValidateEmail() bool {
	regex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return regex.MatchString(r.Email)
}

func (r *Register) ValidatePassword() (bool, string) {
    if len(r.Password) < 8 {
        return false, "At least 8 characters"
    }

    hasUpper := false
    hasNumber := false

    for _, c := range r.Password {
        if unicode.IsUpper(c) {
            hasUpper = true
        }
        if unicode.IsNumber(c) {
            hasNumber = true
        }
    }

    if !hasUpper {
        return false, "At least one mayus"
    }
    if !hasNumber {
        return false, "At least one number"
    }

    return true, ""
}

func (r *Register) ValidateUsername() (bool, string) {
    if len(r.Username) < 3 {
        return false, "At least 8 characters"
    }
    if len(r.Username) > 50 {
        return false, "Max 50 characters"
    }

    regex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !regex.MatchString(r.Username) {
        return false, "Only letters, numbers and underscore"
    }

    return true, ""
}