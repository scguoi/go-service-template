package gracefulstop

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"template/internal/impl"
	"time"

	log "github.com/sirupsen/logrus"
)

// WaitExit will block until os signal happened
func WaitExit(c chan os.Signal, exit func(ctx context.Context)) {
	for i := range c {
		switch i {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("receive exit signal ", i.String(), ",exit...")
			for impl.CurrentReqCount.Load() > 0 {
				log.Infoln("wait for all request done current count:", impl.CurrentReqCount.Load())
				time.Sleep(time.Second)
			}
			exit(context.Background())
			os.Exit(0)
		}
	}
}

// NewShutdownSignal new normal Signal channel
func NewShutdownSignal() chan os.Signal {
	c := make(chan os.Signal)
	// SIGHUP: terminal closed
	// SIGINT: Ctrl+C
	// SIGTERM: program exit
	// SIGQUIT: Ctrl+/
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return c
}
