package service

// Handler
type Handler interface {
	handle(message []byte) error
}
