package network

type Session interface {
	Read() ([]byte,error)
	Write(buff []byte) error
}