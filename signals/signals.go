package signals

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

// SetAppInterruption - TODO
func SetAppInterruption(wg *sync.WaitGroup) {
	wg.Add(1)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)

	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)

	// Un-blocking approach
	// select {
	// case sig := <-c:
	// 	fmt.Println("received message", sig)
	// default:
	// 	fmt.Println("no message received")
	// }

	// Block until any signal is received.
	go func() {
		s := <-c
		fmt.Println("Got signal:", s)

		wg.Done()
	}()

}
