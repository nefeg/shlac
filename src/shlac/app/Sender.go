package app

import "net"

type Sender interface{

	Connect(addr net.Addr)
	IsConnected()(bool)
	Disconnect()
	Send(command string) (received []byte)
}