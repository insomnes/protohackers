package handlers

type Handler interface {
	HandleMessage(msg []byte, verbose bool, remote string) ([]byte, error)
}
