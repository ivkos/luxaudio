package effects

import (
	"github.com/ivkos/luxaudio/internal/led"
	"github.com/ivkos/luxaudio/internal/utils"
	"image/color"
	"log"
	"net"
)

type LuxceptionEffect struct {
	ledCount int
	ledData  []byte
	colors   []byte

	host string
	port uint16

	pinger *utils.Pinger
}

func NewLuxceptionEffect(ledCount int, defaultColor color.RGBA, host string, port uint16, pinger *utils.Pinger) Effect {
	e := &LuxceptionEffect{
		ledCount: ledCount,
		ledData:  make([]byte, ledCount*3),
		colors:   make([]byte, ledCount*3),

		host: host,
		port: port,

		pinger: pinger,
	}

	for i := 0; i < ledCount; i++ {
		e.colors[i*3+0] = byte(float64(defaultColor.G))
		e.colors[i*3+1] = byte(float64(defaultColor.R))
		e.colors[i*3+2] = byte(float64(defaultColor.B))
	}

	go e.listen()

	return e
}

func (e *LuxceptionEffect) Apply(intensities []float64) []byte {
	for i, x := range intensities {
		e.ledData[i*3+0] = byte(float64(e.colors[i*3+0]) * x)
		e.ledData[i*3+1] = byte(float64(e.colors[i*3+1]) * x)
		e.ledData[i*3+2] = byte(float64(e.colors[i*3+2]) * x)
	}

	return e.ledData
}

func (e *LuxceptionEffect) listen() {
	addr, err := utils.GetUDPAddr(e.host, e.port)
	utils.CheckErr(err)

	conn, err := net.ListenUDP("udp", addr)
	utils.CheckErr(err)

	defer func() { _ = conn.Close() }()

	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		go e.handleData(conn, addr, buf[:n])
	}
}

func (e *LuxceptionEffect) handleData(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	if len(data) < 3 {
		log.Printf("Message has invalid length: %d\n", len(data))
		return
	}

	if data[0] != 0x4C || data[1] != 0x58 {
		log.Printf("Invalid header")
		return
	}

	mode := data[2]
	effectPayloadOffset := 3

	if mode == byte(led.Raw) {
		ledCountInPayload := int(data[effectPayloadOffset])
		if ledCountInPayload != e.ledCount {
			log.Printf("Expected %d LEDs, got %d\n", e.ledCount, ledCountInPayload)
			return
		}

		expectedDataLen := 1 + effectPayloadOffset + ledCountInPayload*3
		if len(data) != expectedDataLen {
			log.Printf("Expected %d bytes, got %d\n", expectedDataLen, len(data))
			return
		}

		e.colors = data[effectPayloadOffset+1:]
	} else if mode == byte(led.Ping) {
		var response string
		if e.pinger.IsReachable {
			response = "1"
		} else {
			response = "0"
		}

		_, _ = conn.WriteToUDP([]byte(response), addr)
	} else {
		log.Printf("Unsupported mode: %d\n", mode)
	}
}
