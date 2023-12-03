package fetch

type Fetcher interface {
	GetCh() <-chan []byte
}
