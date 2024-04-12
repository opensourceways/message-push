package hanlder

// Handler
type Handler interface {
	handle(message []byte) error
}
