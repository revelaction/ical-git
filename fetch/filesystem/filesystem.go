package filesystem

import (
	"github.com/revelaction/ical-git/fetch"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem implements the fetch.Fetcher interface
type FileSystem struct {
	rootDir string
	ch      chan fetch.File
}

// New creates a new instance of FileSystem
func New(rootDir string) fetch.Fetcher {
	return &FileSystem{
		rootDir: rootDir,
		ch:      make(chan fetch.File),
	}
}

// GetCh implements the fetch.Fetcher interface method
func (fs *FileSystem) GetCh() <-chan fetch.File {
	go func() {
		defer close(fs.ch)
		filepath.Walk(fs.rootDir, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				fs.ch <- fetch.File{Path: path, Error: err}
				return err
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToLower(info.Name()), ".ical") || strings.HasSuffix(strings.ToLower(info.Name()), ".ics") {
				content, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				fs.ch <- fetch.File{Path: path, Content: content}
			}

			return nil
		})

	}()

	return fs.ch
}
