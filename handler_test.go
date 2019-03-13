package dashbutton

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var testConn net.Conn

func TestMain(m *testing.M) {
	var err error
	testConn, err = net.Dial("udp4", "255.255.255.255:67")
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Println("Closing test connection")
		panic(testConn.Close())
	}()

	os.Exit(m.Run())
}

func TestHandler_AddButtonHandler(t *testing.T) {
	errors := make(chan error)
	buttonChan := make(ButtonHandler)

	handler := NewHandler()
	handler.AddButtonHandler("de:ad:be:ef:ca:fe", buttonChan)

	defer handler.Close()

	go func() {
		if err := handler.Listen(); err != nil {
			log.Printf("Error with handler listen: %v", err)
			t.Fatal(err)
		}
	}()

	go func() {
		select {
		case button := <-buttonChan:
			log.Printf("Received registered button press: %s", button)
			errors <- nil
		case <-time.After(5000 * time.Millisecond):
			errors <- fmt.Errorf("timeout waiting for message")
		}
	}()

	go func() {
		time.Sleep(250 * time.Millisecond)
		if err := sendDHCPPacket("de:ad:be:ef:ca:fe"); err != nil {
			t.Error(err)
		}
	}()

	if err := <-errors; err != nil {
		log.Printf("Received error from channel: %v\n", err)
		t.Error(err)
	}
}

func TestHandler_SetUnknownButtonHandler(t *testing.T) {
	errors := make(chan error)
	unknownButtonChan := make(ButtonHandler)

	handler := NewHandler()
	handler.SetUnknownButtonHandler(unknownButtonChan)

	defer handler.Close()

	go func() {
		if err := handler.Listen(); err != nil {
			log.Printf("Error with handler listen: %v", err)
			t.Fatal(err)
		}
	}()

	go func() {
		select {
		case button := <-unknownButtonChan:
			log.Printf("Received unknown button press: %s", button)
			errors <- nil
		case <-time.After(5000 * time.Millisecond):
			errors <- fmt.Errorf("timeout waiting for message")
		}
	}()

	go func() {
		time.Sleep(250 * time.Millisecond)
		if err := sendDHCPPacket("de:ad:be:ef:ca:fe"); err != nil {
			t.Error(err)
		}
	}()

	if err := <-errors; err != nil {
		t.Error(err)
	}
}

func sendDHCPPacket(hwAddr string) error {
	mac, err := net.ParseMAC(hwAddr)
	if err != nil {
		return err
	}
	dhcp := &layers.DHCPv4{Operation: layers.DHCPOpRequest, HardwareType: layers.LinkTypeEthernet, Xid: 0x12345678,
		ClientIP: net.IP{0, 0, 0, 0}, YourClientIP: net.IP{0, 0, 0, 0}, NextServerIP: net.IP{0, 0, 0, 0}, RelayAgentIP: net.IP{0, 0, 0, 0},
		HardwareLen: 6, ClientHWAddr: mac, ServerName: make([]byte, 64), File: make([]byte, 128),
		Options: []layers.DHCPOption{layers.NewDHCPOption(layers.DHCPOptMessageType, []byte{byte(layers.DHCPMsgTypeDiscover)})}}

	buf := gopacket.NewSerializeBuffer()
	serializeOpts := gopacket.SerializeOptions{}
	if err := gopacket.SerializeLayers(buf, serializeOpts, dhcp); err != nil {
		return err
	}

	log.Printf("Sending DHCP req packet with MAC: %s\n", hwAddr)
	n, err := testConn.Write(buf.Bytes())
	_ = n
	return err
}
