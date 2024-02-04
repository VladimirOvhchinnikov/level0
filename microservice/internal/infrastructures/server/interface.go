package server

type IRouter interface {
	StartServer(path string) error
	Shutdown() error
}
