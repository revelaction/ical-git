package filesystem

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/revelaction/ical-git/fetch"
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
func (fs *FileSystem) GetCh(ctx context.Context) <-chan fetch.File {
	go func() {
		defer close(fs.ch)

		walkErr := filepath.Walk(fs.rootDir, func(path string, info os.FileInfo, err error) error {
			ctxErr := ctx.Err()
			if ctxErr != nil {
				return ctxErr
			}

			if err != nil {
				fs.ch <- fetch.File{Path: path, Error: err}
				return err
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToLower(info.Name()), ".ical") || strings.HasSuffix(strings.ToLower(info.Name()), ".ics") {
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					slog.Error("failed to read file", "path", path, "error", readErr)
					return nil
				}

				fs.ch <- fetch.File{Path: path, Content: content}
			}

			return nil
		})

		if walkErr != nil {
			slog.Error("error walking filesystem", "error", walkErr)
		}
	}()

	return fs.ch
}
