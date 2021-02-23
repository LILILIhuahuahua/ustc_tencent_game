package event

type Session interface {
	Read() ([]byte, error)
	Write(buff []byte) error
}
