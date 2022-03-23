package fileutils

import (
	"embed"
	"io/fs"
	"os"
	"path"
)

func Save(outDir, filePath string, fs embed.FS) error {
	data, err := fs.ReadFile(filePath)
	if err != nil {
		return err
	}

	out, err := CreateLocal(path.Join(outDir, filePath))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func CreateLocal(p string) (*os.File, error) {
	dirPath := path.Dir(p)
	if err := os.MkdirAll(dirPath, fs.ModeDir); err != nil {
		return nil, err
	}
	return os.Create(p)
}
