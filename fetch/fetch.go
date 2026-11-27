package fetch

import "context"

type File struct {
	Path    string
	Content []byte
	Error   error
}

type Fetcher interface {
	GetCh(ctx context.Context) <-chan File
}
