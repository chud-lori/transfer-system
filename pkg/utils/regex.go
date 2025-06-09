package utils

import "regexp"

func ValidateDecimalFormat(input string) bool {
	// This regex matches numbers like "123", "123.45", "0.00"
	// It requires at least one digit before the decimal if a decimal is present.
	// It doesn't allow "123." or ".45"
	regex := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return regex.MatchString(input)
}
