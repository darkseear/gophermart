package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func ValidLuhn(number string) bool {
	sum := 0
	for i := range len(number) {
		digit := int(number[i] - '0')
		if i%2 == len(number)%2 {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		sum += digit
	}
	return sum%10 == 0
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
