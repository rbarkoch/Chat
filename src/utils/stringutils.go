package utils

// StringOrDefault returns the default value if the string is empty.
func StringOrDefault(currentValue string, defaultValue string) string {
	if currentValue == "" {
		return defaultValue
	}

	return currentValue
}
