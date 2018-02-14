package socket

import (
	"net"
	"io"
	"log"
	"os"
	capi "shlacd/cli"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"github.com/urfave/cli"
	"errors"
	"github.com/mattn/go-shellwords"
	"regexp"
)

type handler struct {
	addr    net.Addr
}

const WlcMessage = "ShLAC terminal connected OK\n" +
	"type \"help\" or \"\\h\" for show available commands"
const logPrefix = "[client.telnet]"

var ErrConnectionClosed = errors.New("** command <QUIT> received")


func NewHandler(listen net.Addr) *handler{

	return &handler{ addr:listen }
}


func (h *handler) Handle(ctx capi.Context){

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

				h.handleConnection(Connection, ctx)
				Connection.Close()

				log.Printf(logPrefix + "Client connection closed [connid:%v]", Connection)
			}()

		}else{
			log.Println(logPrefix, err.Error())
			continue
		}
	}
}

func (h *handler)handleConnection(Connection net.Conn, ctx capi.Context){

	var response string
	var Responder = NewClient(Connection)

	defer func(response *string){

		if r := recover(); r != nil{

			if r == io.EOF {
				*response = "client socket closed."
				Responder.WriteString("\n" + (*response) + "\n")
				slog.InfoLn(logPrefix + "Session closed by cause: " + (*response))

			}else{
				slog.InfoLn(logPrefix + "Session closed by cause: " , r)
			}
		}else{
			Responder.WriteString("\n" + (*response) + "\n")
			slog.InfoLn(logPrefix + "Session closed by cause: " + (*response))
		}
	}(&response)


	Responder.WriteString(WlcMessage)

	Cli := capi.New()

	Cli.Writer              = Responder
	Cli.ErrWriter           = Responder

	// COMMANDS
	Cli.Commands = []cli.Command{
		capi.NewComAdd(&ctx),
		capi.NewComExport(&ctx),
		capi.NewComRemove(&ctx),
		capi.NewComPurge(&ctx),
		capi.NewComGet(&ctx),
		{
			Name:    "exit",
			Aliases: []string{`q`},
			Usage:   "close connection",
			UsageText: "Example: " ,

			Action:  func(c *cli.Context) error {

				slog.DebugLn("Action: exit")

				c.App.Writer.Write([]byte("Sending <QUIT> signal..."))
				panic(ErrConnectionClosed)

				return nil
			},
		},
	}

	Cli.After = func(c *cli.Context) error {

		c.App.Writer.Write(PacketTerm)

		return nil
	}

	Cli.ExitErrHandler = func(c *cli.Context, err error){
		c.App.Writer.Write([]byte(err.Error()))
		slog.DebugLn(logPrefix, err)
	}



	for{

		if rcb, err := Responder.ReadData(); len(rcb) != 0{

			if err != nil {
				slog.CritLn(err.Error())
				response = err.Error()
				Responder.WriteString(response)

			}else{

				slog.DebugLn(logPrefix, "Args (byte,raw):", rcb)
				slog.DebugLn(logPrefix, "Args (string,raw):", string(rcb))

				if match,_ := regexp.Match(`^\w.*`, rcb); match != true{
					rcb = []byte("help")
					slog.DebugLn(logPrefix, "Incorrect args, show help")
				}

				args,_ := shellwords.Parse("self " + string(rcb))

				Cli.Run( args )

				slog.DebugLn(logPrefix, "Cli.Run (complete)")
			}
		}
	}
}



