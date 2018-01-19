package socket

import (
	"net"
	"io"
	"log"
	"os"
	ComLineIf "shlacd/cli"
	"shlacd/cli/Com"
	"shlacd/app/api"
)

type handler struct {

	addr    net.Addr
	cli     ComLineIf.CLI
}

const WlcMessage = "ShLAC terminal connected OK\n" +
	"type \"help\" or \"\\h\" for show available commands"
const logPrefix = "[client.telnet]"


func NewHandler(listen net.Addr, cli ComLineIf.CLI) *handler{

	return &handler{ addr:listen, cli:cli }
}


func (h *handler) Handle(Tab api.TimeTable){

	IPC, err := net.Listen(h.addr.Network(), h.addr.String())
	if err != nil {
		log.Panicf("%s: %s", "ERROR", err.Error())
	}else{
		log.Printf(logPrefix + "Listen: %s://%s", IPC.Addr().Network(), IPC.Addr().String())
	}
	defer func(){
		IPC.Close()
		if UAddr, err := net.ResolveUnixAddr(h.addr.Network(), h.addr.String()); err == nil{
			os.Remove(UAddr.String())
		}
	}()

	for{
		if Connection, err := IPC.Accept(); err == nil {

			go func(){
				log.Printf(logPrefix + "New client connection accepted [connid:%v]", Connection)

				h.handleConnection(Connection, Tab)
				Connection.Close()

				log.Printf(logPrefix + "Client connection closed [connid:%v]", Connection)
			}()

		}else{
			log.Println(err.Error())
			continue
		}
	}
}

func (h *handler)handleConnection(Connection net.Conn, Tab api.TimeTable){

	var response string

	defer func(response *string){

		if r := recover(); r != nil{

			if r == io.EOF {
				*response = "client socket closed."
				writeData(Connection, "\n"+(*response)+"\n")
				log.Println(logPrefix + "Session closed by cause: " + (*response))

			}else{
				log.Println(logPrefix + "Session closed by cause: " , r)
			}
		}else{
			writeData(Connection, "\n" + (*response) + "\n")
			log.Println(logPrefix + "Session closed by cause: " + (*response))
		}
	}(&response)


	writeData(Connection, WlcMessage)
	for{

		if rcv, err := readData(Connection); rcv != ""{

			if err != nil {
				log.Println(err.Error())
				response = err.Error()

			}else{

				response = "Unknown command"
				if Command, args, err := h.cli.Resolve(rcv); Command != nil{

					response, err = Command.Exec(Tab, args)

					if err != nil{
						response = err.Error()

						// ComQuit
						if err == Com.ErrConnectionClosed{
							return
						}
					}
				}else{
					response += "\n" + h.cli.Help()
				}
			}

			writeData(Connection, response)
		}
	}
}



