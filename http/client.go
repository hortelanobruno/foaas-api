package http

type Client interface {
	Get(url string) ([]byte, error)
}
