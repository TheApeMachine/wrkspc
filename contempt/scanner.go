package contempt

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/machine"
	"golang.org/x/sync/semaphore"
)

/*
Scanner is an object that can scan the network for live machines that we can try to connect
to and automatically provision.
*/
type Scanner struct {
	ipSeq func() chan string
	lock  *semaphore.Weighted
}

/*
NewScanner constructs a network Scanner and returns a reference to it.
*/
func NewScanner(ipRange *Range) *Scanner {
	errnie.Traces()

	return &Scanner{
		ipSeq: NewSequencer(ipRange, "192.168.1.").Generate,
		lock:  semaphore.NewWeighted(machine.NewSystem().Ulimit),
	}
}

/*
Sweep the network using various scanning methods to discover connectable
interfaces so we can provision the machine as cluster fodder.
*/
func (scanner *Scanner) Sweep() chan Connection {
	errnie.Traces()
	out := make(chan Connection)

	go func() {
		defer close(out)
		var wg sync.WaitGroup

		for ip := range scanner.ipSeq() {
			scanner.lock.Acquire(context.TODO(), 1)

			wg.Add(1)
			go func(ip string, wg *sync.WaitGroup) {
				defer wg.Done()

				target := fmt.Sprintf("%s:%d", ip, 22)
				conn, err := net.DialTimeout("tcp", target, 1*time.Second)

				// Connection error, bail!
				if !errnie.Handles(err).With(errnie.NOLO).OK {
					return
				}

				conn.Close()
				scanner.lock.Release(1)
				errnie.Logs(target).With(errnie.INFO)

				// Construct a new Connection and send it to the caller.
				out <- NewConnection(SSH{IP: ip})
			}(ip, &wg)
		}

		wg.Wait()
	}()

	return out
}
