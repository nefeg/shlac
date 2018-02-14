package sig

import (
	"os"
	"syscall"
	"os/signal"
	"errors"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

var ErrSigINT = errors.New("SIG_INT")

func SIG_INT(callback *func()){

	go func() {

		// slog.Debugln("[shared.sig] Listening SIGINT")

		sig := make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGINT)

		<-sig

		slog.Infoln("[shared.sig] Received termination signal:", ErrSigINT)
		slog.Infoln("[shared.sig] Waiting for signal handler...")

		(*callback)()

		slog.Infoln("[shared.sig] Exit (SIGINT)")

		os.Exit(0)
	}()
}