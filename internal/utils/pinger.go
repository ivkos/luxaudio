package utils

import (
	"github.com/ivkos/luxaudio/internal/led"
	"log"
	"net"
	"time"
)

type Pinger struct {
	conn     *net.UDPConn
	interval time.Duration
	timeout  time.Duration
	verbose  bool

	IsReachable bool
}

func NewPinger(conn *net.UDPConn, interval time.Duration, verbose bool) *Pinger {
	pinger := &Pinger{
		conn:     conn,
		interval: interval,
		timeout:  1 * time.Second,
		verbose:  verbose,

		IsReachable: true,
	}

	go pinger.start()

	return pinger
}

func (pinger *Pinger) start() {
	timer := time.NewTimer(0)
	pingPayload := led.MakePingPayload()

	for {
		timer.Reset(pinger.interval)
		<-timer.C

		_, err := pinger.conn.Write(pingPayload)
		if err != nil {
			pinger.setReachable(false)
			pinger.logVerbose("WARN: Could not write ping payload: %v", err)
			continue
		}

		err = pinger.conn.SetReadDeadline(time.Now().Add(pinger.timeout))
		if err != nil {
			pinger.setReachable(false)
			pinger.logVerbose("WARN: Could not set ping deadline: %v", err)
			continue
		}

		result := make([]byte, 1)
		n, err := pinger.conn.Read(result)
		if err != nil {
			pinger.logVerbose("WARN: Could not read ping response: %v", err)
			pinger.setReachable(false)
			continue
		}

		if n != 1 {
			pinger.logVerbose("WARN: Ping response has unexpected length %d", n)
			pinger.setReachable(false)
			continue
		}

		if result[0] != '1' {
			pinger.logVerbose("WARN: Ping response is unexpected: %x", result[0])
			pinger.setReachable(false)
			continue
		}

		pinger.setReachable(true)
	}
}

func (pinger *Pinger) setReachable(reachable bool) {
	if reachable != pinger.IsReachable {
		log.Printf("Reachable = %t", reachable)
	}

	pinger.IsReachable = reachable
}

func (pinger *Pinger) logVerbose(format string, v ...interface{}) {
	if pinger.verbose {
		log.Printf(format, v...)
	}
}
