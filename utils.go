package tvchooser

import (
	"path/filepath"
)

// splitPathAndName separates a string into a path and a name based on the last occurrence of a separator rune.
// For example, if the string is "D:\\user\\Music|Música" and the separator is '|', it will return "D:\\user\\Music" and "Música".
// If the separator is not found, it returns the original string as the path and an the base path name (filepath.Base(path)) as the name.
func splitPathAndName(pathAndName string, separator rune) (path string, name string) {
	// Find the last occurrence of the separator rune in the string.
	lastSeparatorIndex := -1
	for i, r := range pathAndName {
		if r == separator {
			lastSeparatorIndex = i
		}
	}

	// If the separator is not found, return the original string as the path and the base name as the name.
	if lastSeparatorIndex == -1 {
		return pathAndName, filepath.Base(pathAndName)
	}

	// Split the string into path and name based on the last occurrence of the separator.
	path = pathAndName[:lastSeparatorIndex]
	name = pathAndName[lastSeparatorIndex+1:]

	return path, name
}
