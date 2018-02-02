package socket

import (
	"net"
	"strings"
	"log"
)

func readData(Connection net.Conn) (rcv string, err error){

	log.Println("[SYS]readData: ....")


	tmp := make([]byte, 4096)

	length, err := Connection.Read(tmp[:])
	if err != nil {
		panic(err)
	}

	if length>0{
		rcv = strings.TrimSpace(string(tmp[:length]))
	}

	log.Println("[SYS]readData (string):", rcv)
	log.Println("[SYS]readData (raw):", []byte(rcv))
	log.Println("[SYS]readData (error):", err)

	return rcv, err
}

func writeData(Connection net.Conn, data string) (int, error){

	log.Println("[SYS]writeData (string):", data)

	response := append([]byte(data), []byte{00, 10, 62, 62}...)

	log.Println("[SYS]writeData (raw+term):", response)


	return Connection.Write(response)
}
