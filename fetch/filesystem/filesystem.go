package filesystem

import (
	"os"
	"path/filepath"
	"github.com/revelaction/ical-git/fetch"
)

// FileSystem implements the fetch.Fetcher interface
type FileSystem struct {
	rootDir string
	ch      chan []byte
}

// New creates a new instance of FileSystem
func New(rootDir string) *FileSystem {
	return &FileSystem{
		rootDir: rootDir,
		ch:      make(chan []byte),
	}
}

// GetCh implements the fetch.Fetcher interface method
func (fs *FileSystem) GetCh() <-chan []byte {
	go func() {
        // TODO 
		filepath.Walk(fs.rootDir, func(path string, _ os.FileInfo, _ error) (err error) {
			fs.ch <- path
			return
		})

		defer close(fs.ch)
	}()

	return fs.ch
}
