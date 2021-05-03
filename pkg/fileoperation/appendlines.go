package fileoperation

import (
	"os"
)

// AppendLines appends all lines to a text file, creating a new file if it doesn't exist.
func AppendLines(lines []string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return writelines(lines, file, 64*1024)
}
