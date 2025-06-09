package validator

import "regexp"

func ValidateDecimalFormat(input string) bool {
	// assuming that the input will always in number with 5 digits after the decimal point
	// e.g. 12345.12345
	regex := regexp.MustCompile(`^\d+(\.\d{5})?$`)
	return regex.MatchString(input)
}
