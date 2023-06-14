package app

import (
	"strconv"
	"strings"
)

func IsInteger(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

// check if str is uuid with minimum cost
func IsUUID(str string) bool {
	if len(str) != 36 {
		return false
	}

	if str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		return false
	}

	validChars := "0123456789abcdefABCDEF"
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !strings.ContainsRune(validChars, rune(str[i])) {
			return false
		}
	}

	return true
}
