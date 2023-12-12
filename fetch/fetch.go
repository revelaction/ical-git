package fetch

type File struct {
	Path string
    Content []byte
    Error error
}

type Fetcher interface {
	GetCh() <-chan File
}


