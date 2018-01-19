package lib

import (
	"net"
	"os"
)


type libCommand struct{
	connection  net.Conn
}

func NewCommander() *libCommand{

	return &libCommand{}
}

func (l *libCommand) Send(command string) ([]byte){

	if !l.IsConnected() {
		panic("connection not exist")
	}

	l.read() // clear socket buffer

	return l.write([]byte(command+"\n")).read()
}

//func (l *libCommand) SendOnce(command string) ([]byte){
//
//	if !l.IsConnected() {
//		l.Connect()
//		defer func(){l.Disconnect()}()
//	}
//
//	l.read() // clear socket buffer
//
//	l.write([]byte(command+"\n")).read()
//
//	return l.write([]byte(command+"\n")).read()
//}


func (l *libCommand) Connect(addr net.Addr) {

	if l.IsConnected(){
		panic("connection already exist")
	}

	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		panic(err)
	}

	l.connection = conn
}

func (l *libCommand) IsConnected() bool{

	return l.connection != nil
}

func (l *libCommand) Disconnect(){
	l.connection.Write([]byte(`\q`))

	l.connection.Close()

	// remove socket-file if connected via unix-socket
	if UAddr, err := net.ResolveUnixAddr(l.connection.LocalAddr().Network(), l.connection.LocalAddr().String()); err == nil{
		os.Remove(UAddr.String())
	}

}


func (l *libCommand) read() (flushed []byte){

	bufSize := 256
	buf := make([]byte, bufSize)

	for{
		n,e := l.connection.Read(buf)

		flushed = append(flushed, buf[:n]...)

		if e != nil || n < bufSize {break}
	}

	return flushed
}

func (l *libCommand) write(data []byte) *libCommand{

	l.connection.Write(data)

	return l
}