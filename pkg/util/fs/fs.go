package fs

import (
	"io/fs"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func OpenFile(path string) (fs.File, error) {
	return os.OpenFile(path, os.O_RDONLY, 0644)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func RemoveFile(path string) {
	_ = os.Remove(path)
}

func WriteFile(path string, data []byte, m fs.FileMode) error {
	return os.WriteFile(path, data, m)
}

func MustWriteFile(path string, data []byte, m fs.FileMode) {
	if err := os.WriteFile(path, data, m); err != nil {
		panic(err)
	}
}

func ReadDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}

func MkDir(path string) error {
	return os.MkdirAll(path, 0755)
}
