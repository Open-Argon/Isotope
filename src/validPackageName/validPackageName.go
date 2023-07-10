package validpackagename

import "strings"

func ValidPackageName(str string) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", "-")
	str = strings.ReplaceAll(str, "_", "-")
	str = strings.ReplaceAll(str, "/", "-")
	str = strings.ReplaceAll(str, "\\", "-")

	// Remove all non-alphanumeric characters
	for i := 0; i < len(str); i++ {
		if !((str[i] >= 'a' && str[i] <= 'z') || (str[i] >= '0' && str[i] <= '9') || str[i] == '-' || str[i] == '.') {
			str = str[:i] + str[i+1:]
			i--
		}
	}
	return str
}
