package validators

import (
	"unicode"
)

const DefaultPasswordLength int = 8

func VerifyPassword(s string) bool {
	chars := 0
	var (
		properLength bool
		number       bool
		letter       bool
		upper        bool
	)
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLetter(c) || c == ' ':
			letter = true
		default:
		}
		chars++
	}
	properLength = chars >= DefaultPasswordLength
	return number && letter && upper && properLength
}
