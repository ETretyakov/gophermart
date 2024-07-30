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
		upper        bool
	)
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			chars++
		case unicode.IsUpper(c):
			upper = true
			chars++
		case unicode.IsLetter(c) || c == ' ':
			chars++
		default:
			return false
		}
	}
	properLength = chars >= DefaultPasswordLength
	return number && upper && properLength
}
