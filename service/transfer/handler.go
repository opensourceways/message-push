package transfer

// Handler
type Handler interface {
	handle(message []byte) error
}
