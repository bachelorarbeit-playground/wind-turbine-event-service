package file

import "os"

func CheckIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
