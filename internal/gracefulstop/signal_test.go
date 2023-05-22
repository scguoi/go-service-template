// BEGIN: 8f6a4b3c4d5e
package gracefulstop_test

import (
	"context"
	"sync"
	"syscall"
	"template/internal/gracefulstop"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestWaitExit(t *testing.T) {
	// create a shutdown signal channel
	c := gracefulstop.NewShutdownSignal()

	// create a wait group to wait for the exit function to complete
	var wg sync.WaitGroup
	wg.Add(1)

	// create an exit function that will be called when a shutdown signal is received
	exit := func(ctx context.Context) {
		log.Println("exit function called")
		wg.Done()
	}

	// start the wait exit function in a separate goroutine
	go gracefulstop.WaitExit(c, exit)

	// wait for a short time to allow the goroutine to start
	time.Sleep(100 * time.Millisecond)

	// send a SIGINT signal to the process
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	// wait for the exit function to complete
	wg.Wait()
}

func TestNewShutdownSignal(t *testing.T) {
	// create a shutdown signal channel
	c := gracefulstop.NewShutdownSignal()

	// send a SIGINT signal to the process
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	// wait for the signal to be received
	select {
	case <-c:
		// signal received, test passed
	case <-time.After(1 * time.Second):
		// signal not received, test failed
		t.Error("shutdown signal not received")
	}
}

// END: 8f6a4b3c4d5e
