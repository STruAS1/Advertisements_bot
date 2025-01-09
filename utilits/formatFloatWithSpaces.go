package utilits

import (
	"fmt"
	"strings"
)

func FormatFloatWithSpaces(number float64) string {
	numStr := fmt.Sprintf("%.2f", number)

	parts := strings.Split(numStr, ".")
	intPart := parts[0]
	fracPart := parts[1]

	var result strings.Builder

	for i, digit := range intPart {
		if (len(intPart)-i)%3 == 0 && i != 0 {
			result.WriteRune(' ')
		}
		result.WriteRune(digit)
	}

	if fracPart != "" {
		result.WriteString("." + fracPart)
	}

	return result.String()
}
