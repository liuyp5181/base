package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var closeFunc []func()

func init() {
	go func() {
		osc := make(chan os.Signal, 1)
		signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
		<-osc
		for _, f := range closeFunc {
			f()
		}
		os.Exit(1500)
	}()
}

func RegisterClose(f func()) {
	closeFunc = append(closeFunc, f)
}
