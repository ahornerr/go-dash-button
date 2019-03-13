package main

import (
	"log"
	"net"

	"github.com/krolaw/dhcp4"
)

// ButtonHandlerFunc defines the function which will be called when a button is pressed
type ButtonHandlerFunc func(net.HardwareAddr)

// DashButtonHandler listens for DHCP requests coming from an Amazon Dash button. The handler provides
type DashButtonHandler struct {
	buttonHandlers map[string]ButtonHandlerFunc
	unknownHandler ButtonHandlerFunc
}

// DefaultUnknownButtonHandler logs that an unknown dash button was pressed, along with it's MAC address
func DefaultUnknownButtonHandler(hwAddr net.HardwareAddr) {
	log.Printf("Unknown dash button pressed. MAC address: %s", hwAddr)
	// Output: Unknown dash button pressed. MAC address: [MAC address]
}

// NewDashButtonHandler creates a new instance of the button handler
func NewDashButtonHandler() DashButtonHandler {
	return DashButtonHandler{
		buttonHandlers: map[string]ButtonHandlerFunc{},
		unknownHandler: DefaultUnknownButtonHandler,
	}
}

// AddButtonHandler adds a ButtonHandlerFunc to listen for a buttonMAC
func (dbh DashButtonHandler) AddButtonHandler(buttonMAC string, handlerFunc ButtonHandlerFunc) {
	dbh.buttonHandlers[buttonMAC] = handlerFunc
}

// SetUnknownButtonHandler sets the handler function to be called when an unknown button is pressed
func (dbh DashButtonHandler) SetUnknownButtonHandler(handlerFunc ButtonHandlerFunc) {
	dbh.unknownHandler = handlerFunc
}

func (dbh DashButtonHandler) ServeDHCP(req dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) dhcp4.Packet {
	// For debugging purposes
	// fmt.Println(gopacket.NewPacket(req, layers.LayerTypeDHCPv4, gopacket.Default).String())

	hwAddr := req.CHAddr()
	if handler := dbh.buttonHandlers[hwAddr.String()]; handler != nil {
		handler(hwAddr)
		return nil
	}
	if dbh.unknownHandler != nil {
		dbh.unknownHandler(hwAddr)
	}
	return nil
}
