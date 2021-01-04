package network

type Conn interface {
     Read() ([]byte,error)
     Write([] byte)
}