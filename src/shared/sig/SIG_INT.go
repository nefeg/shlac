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

		// slog.DebugLn("[shared.sig] Listening SIGINT")

		sig := make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGINT)

		<-sig

		slog.InfoLn("[shared.sig] Received termination signal:", ErrSigINT)
		slog.InfoLn("[shared.sig] Waiting for signal handler...")

		(*callback)()

		slog.InfoLn("[shared.sig] Exit (SIGINT)")

		os.Exit(0)
	}()
}