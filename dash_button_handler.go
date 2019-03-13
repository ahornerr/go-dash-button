package main

import (
	"log"
	"net"

	"github.com/krolaw/dhcp4"
)

type ButtonHandlerFunc func(net.HardwareAddr)

type dashButtonHandler struct {
	buttonHandlers map[string]ButtonHandlerFunc
	unknownHandler ButtonHandlerFunc
}

func defaultUnknownButtonHandler(hwAddr net.HardwareAddr) {
	log.Printf("Unknown dash button: %s", hwAddr)
}

func NewDashButtonHandler() *dashButtonHandler {
	return &dashButtonHandler{
		buttonHandlers: map[string]ButtonHandlerFunc{},
		unknownHandler: defaultUnknownButtonHandler,
	}
}

func (dbh dashButtonHandler) AddButtonHandler(buttonMAC string, handlerFunc ButtonHandlerFunc) {
	dbh.buttonHandlers[buttonMAC] = handlerFunc
}

func (dbh dashButtonHandler) SetUnknownButtonHandler(handlerFunc ButtonHandlerFunc) {
	dbh.unknownHandler = handlerFunc
}

func (dbh dashButtonHandler) ServeDHCP(req dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) dhcp4.Packet {
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
