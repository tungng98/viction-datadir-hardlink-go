package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstDir := filepath.Dir(dstPath)
	if !IsExist(dstDir) {
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func GetAbsPath(fPath string) (string, error) {
	absolutePath, err := filepath.Abs(fPath)
	if err == nil {
		absolutePath = NormalizePath(absolutePath)
	}
	return absolutePath, err
}

func IsExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return !os.IsNotExist(err)
}

func IsFileExist(fPath string) bool {
	fileInfo, err := os.Stat(fPath)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

func NormalizePath(fPath string) string {
	newPath := strings.ReplaceAll(fPath, "\\", "/") // enfore linux path style for clarity
	return newPath
}

func WriteFile(filePath string, data []byte) error {
	dir := filepath.Dir(filePath)
	if !IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, data, 0644)
}
