package signalhandler

import (
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	// Create a channel, that when closed signals
	// other stuff they need to stop
	var stopChan = make(chan struct{})

	// Create a closure, literally.
	f := func(os.Signal) { close(stopChan) }

	// Install our closure as a handler for when we feel a Ctrl-C, etc.
	done := Simple(f)
	defer done() // eventually let the signal handler know it should exit.

	// continue on our lives, checking for <-stopChan where necessary

}

func Test_Signals(t *testing.T) {
	defer leaktest.Check(t)()

	var (
		count atomic.Int64
		h     = func(os.Signal) { count.Add(1) }
	)

	// Get a handle on our process
	us, _ := os.FindProcess(os.Getpid())

	Convey("When a signalhandler is installed", t, func() {
		done := Simple(h)
		defer done()

		Convey("and the signal is tripped, everything is normal", func() {
			So(count.Load(), ShouldEqual, 0)
			So(us.Signal(syscall.SIGINT), ShouldBeNil)
			time.Sleep(10 * time.Millisecond) // pause
			So(count.Load(), ShouldEqual, 1)
		})
	})

	Convey("When a signalhandler is installed", t, func() {
		done := Simple(h)
		defer done()

		Convey("and time passes, eventually exiting without the signal being tripped, and everything is normal.", func() {
			time.Sleep(10 * time.Millisecond) // pause
		})
	})
}
