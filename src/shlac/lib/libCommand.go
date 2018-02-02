package lib

import (
	"net"
	"os"
	"github.com/umbrella-evgeny-nefedkin/slog"
)


type libCommand struct{
	connection  net.Conn
}

func NewCommander() *libCommand{

	return &libCommand{}
}

func (l *libCommand) Send(command string) ([]byte){

	slog.DebugF("[Sender(libCommand)] Send: `%s`\n", command)

	if !l.IsConnected() {
		panic("connection not exist")
	}

	return l.write([]byte(command+"\n")).read()
}


func (l *libCommand) Connect(addr net.Addr) {

	if l.IsConnected(){
		panic("connection already exist")
	}

	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		panic(err)
	}

	l.connection = conn
	l.read() // flush socket
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

	slog.DebugLn("[Sender(libCommand)] read: ....")

	flushed = []byte{}
	bufSize := 256
	buf := make([]byte, bufSize)

	for{
		slog.DebugLn("[Sender(libCommand)] read: *** loop ***")

		n,e := l.connection.Read(buf)

		slog.DebugF("[Sender(libCommand)] read (loop): %d bytes\n", n)
		slog.DebugLn("[Sender(libCommand)] read (loop): error:", e)

		flushed = append(flushed, buf[:n]...)

		if e != nil || n < bufSize {break}
	}

	slog.DebugF("[Sender(libCommand)] read(flushed): %d bytes\n", len(flushed))
	slog.DebugLn("[Sender(libCommand)] read (raw): ", flushed)
	slog.DebugLn("[Sender(libCommand)] read (string): ", string(flushed))

	return flushed
}

func (l *libCommand) write(data []byte) *libCommand{

	slog.DebugLn("[Sender(libCommand)] write (raw): ", data)
	slog.DebugLn("[Sender(libCommand)] write (string): ", string(data))

	n, err := l.connection.Write(data)

	slog.DebugF("[Sender(libCommand)] write: %d bytes\n", n)
	slog.DebugLn("[Sender(libCommand)] write (error): ", err)


	return l
}