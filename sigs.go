// Package signalhandler is used to make it easier to safely and consistently use
// a signal handler.
package signalhandler

import (
	"os"
	"os/signal"
	"syscall"
)

// HandlerFunc is a definition of a closure one can provides to the Handler for when a representative Signal is seen.
type HandlerFunc func(os.Signal)

// DoneFunc is a definition for a closure returned to signal we should stop our listening.
type DoneFunc func()

// Simple installs the HandlerFunc, listening only once for SIGINT or SIGKILL and then stops.
func Simple(f HandlerFunc) DoneFunc {
	return Handler(f, true, syscall.SIGINT, syscall.SIGKILL)
}

// Handler installs the HandlerFunc, optionally stopping when a signal is received, for a list of specified signals.
func Handler(f HandlerFunc, stopOnSignal bool, sigs ...os.Signal) DoneFunc {
	var sigChan = make(chan os.Signal, 1)

	// Stream the signals we care about
	signal.Notify(sigChan, sigs...) // syscall.SIGINT, syscall.SIGTERM

	go func() {
		defer signal.Stop(sigChan) // ensure we disconnect
		for s := range sigChan {
			go f(s) // Fire-and-Forget
			if stopOnSignal {
				return
			}
		}
	}()

	return func() { close(sigChan) }
}
