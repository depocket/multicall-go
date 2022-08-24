package utils

import "strings"

func CleanSpaces(input string) string {
	fields := strings.Fields(input)
	var result string
	for i := 0; i < len(fields); i++ {
		result += fields[i]
	}
	return result
}
