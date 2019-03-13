// Package dashbutton handles listens for DHCP requests and calls registered handlers when notified.
// Typically this will be used for Amazon Dash buttons which, when pressed, connect to WiFI and send a DHCP request.
// Devices are identified by their hardware address (MAC address) and handlers are registered as so.
package dashbutton

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ButtonHandler defines the function which will be called when a button is pressed
type ButtonHandler chan net.HardwareAddr

// Handler listens for DHCP requests coming from an Amazon Dash button
type Handler struct {
	buttonHandlers map[string]ButtonHandler
	unknownHandler ButtonHandler
	quitChan       chan bool
	exitChan       chan bool
}

// NewHandler creates a new instance of the button handler
func NewHandler() Handler {
	return Handler{
		buttonHandlers: map[string]ButtonHandler{},
		quitChan:       make(chan bool),
		exitChan:       make(chan bool),
	}
}

// AddButtonHandler adds a ButtonHandler to listen for a buttonMAC
func (h Handler) AddButtonHandler(buttonMAC string, handlerChan ButtonHandler) {
	h.buttonHandlers[strings.ToLower(buttonMAC)] = handlerChan
}

// SetUnknownButtonHandler sets the handler function to be called when an unknown button is pressed
func (h *Handler) SetUnknownButtonHandler(handlerFunc ButtonHandler) {
	h.unknownHandler = handlerFunc
}

// Listen begins listening for DHCP requests. Call this once button handlers are registered. Blocking.
func (h Handler) Listen() error {
	conn, err := net.ListenPacket("udp4", ":67")
	if err != nil {
		return err
	}
	defer conn.Close()
	buf := make([]byte, 1500)
	for {
		select {
		case <-h.quitChan:
			close(h.exitChan)
			return nil
		default:
			if err := conn.SetDeadline(time.Now().Add(1 + time.Second)); err != nil {
				return err
			}
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				fmt.Println("Failed to accept connection:", err.Error())
				return err
			}
			packet := gopacket.NewPacket(buf[:n], layers.LayerTypeDHCPv4, gopacket.Default)
			if dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4); dhcpLayer != nil {
				// For debugging purposes
				// fmt.Printf("Got DHCP packet: %s", packet.String())
				hwAddr := dhcpLayer.(*layers.DHCPv4).ClientHWAddr
				if handler := h.buttonHandlers[hwAddr.String()]; handler != nil {
					select {
					case handler <- hwAddr:
					default:
						log.Println("Dropping write to registered button handler (channel not receiving)")
					}
				} else if h.unknownHandler != nil {
					select {
					case h.unknownHandler <- hwAddr:
					default:
						log.Println("Dropping write to unknown button handler (channel not receiving)")
					}
				} else {
					log.Printf("No button handler for MAC: %s\n", hwAddr)
				}
			}
		}
	}
}

// Close shuts down the DHCP server and stops listening. Blocks until the listener is actually stopped.
func (h Handler) Close() {
	select {
	case <-h.quitChan:
		log.Println("Close quit already called")
	default:
		close(h.quitChan)
		<-h.exitChan // Block until actually exited
	}
}
