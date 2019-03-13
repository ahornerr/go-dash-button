package main

import (
	"fmt"
	"log"
	"net"

	"github.com/krolaw/dhcp4"
)

func main() {
	handler := NewDashButtonHandler()
	handler.AddButtonHandler("fc:65:de:b2:8c:df", func(hwAddr net.HardwareAddr) {
		fmt.Println("Yay! My dash button pressed!")
	})

	fmt.Println("Listening for Dash buttons...")
	log.Fatal(dhcp4.ListenAndServe(handler))
}
