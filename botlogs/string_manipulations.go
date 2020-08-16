package botlogs

import "os"

func isDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	if f.Mode().IsDir() {
		return true
	}
	return false
}
