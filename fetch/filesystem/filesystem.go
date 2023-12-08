package filesystem

import (
	"github.com/revelaction/ical-git/fetch"
	"os"
	"path/filepath"
	"strings"
    "log/slog"
)

// FileSystem implements the fetch.Fetcher interface
type FileSystem struct {
	rootDir string
	ch      chan []byte
}

// New creates a new instance of FileSystem
func New(rootDir string) fetch.Fetcher {
	return &FileSystem{
		rootDir: rootDir,
		ch:      make(chan []byte),
	}
}

// GetCh implements the fetch.Fetcher interface method
func (fs *FileSystem) GetCh() <-chan []byte {
	go func() {
		// TODO
		filepath.Walk(fs.rootDir, func(path string, info os.FileInfo, _ error) (err error) {

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToLower(info.Name()), ".ical") || strings.HasSuffix(strings.ToLower(info.Name()), ".ics") {
				content, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				fs.ch <- content
                slog.Info("🗓️  icalendar file", "path", path)
			}

			return nil
		})

		defer close(fs.ch)
	}()

	return fs.ch
}
