package socket

import "net"

func NewClient(connection net.Conn) *responder{

	return &responder{connection, &[]byte{}, 0}
}


type responder struct{
	connection  net.Conn
	buf         *[]byte
	last        int
}


func (r *responder) Write(b []byte) (n int, err error){

	// slog.DebugLn("[client.telnet] Write :", len(p), r.last)

	r.connection.Write(b)

	return len(b), nil
}



func (r *responder) ReadData() ([]byte, error){

	return readData(r.connection)
}

func (r *responder) ReadString() (string, error){

	b, err := readData(r.connection)

	return string(b), err
}

func (r *responder) WriteData(data []byte) (int, error){

	return writeData(r.connection, data)
}

func (r *responder) WriteString(str string) (int, error){

	return writeData(r.connection, []byte(str))
}
