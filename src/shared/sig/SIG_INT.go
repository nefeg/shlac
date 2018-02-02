package sig

import (
	"os"
	"syscall"
	"os/signal"
	"errors"
	"log"
)

var ErrSigINT = errors.New("SIG_INT")

//noinspection GoSnakeCaseUsage
func SIG_INT(callback func()){

	go func() {

		sig := make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGINT)

		<-sig

		log.Println("Received termination signal:", ErrSigINT)

		if callback != nil{
			callback()
		}

		os.Exit(0)
	}()
}