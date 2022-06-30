package service

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/romik1505/fileserver/internal/app/config"
)

type FileService struct {
	root http.FileSystem
}

type IFileService interface {
	http.FileSystem
	SaveFile(ctx context.Context, in io.Reader, filename string) error
}

func NewFileService() *FileService {
	return &FileService{
		root: http.FileSystem(http.Dir(config.GetValue(config.RootFileDirectory))),
	}
}

func (fsys FileService) Open(path string) (http.File, error) {
	f, err := fsys.root.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, fs.ErrPermission
	}
	return f, nil
}

func (fsys *FileService) SaveFile(ctx context.Context, in io.Reader, filename string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s", fsys.root, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, in); err != nil {
		return err
	}
	return nil
}
