package main

import (
	"fmt"
	"log"

	"github.com/ahornerr/go-dash-button"
)

func main() {
	buttonChan := make(dashbutton.ButtonHandler)
	unknownButtonChan := make(dashbutton.ButtonHandler)

	handler := dashbutton.NewHandler()
	handler.SetUnknownButtonHandler(unknownButtonChan)
	handler.AddButtonHandler("fc:65:de:b2:8c:df", buttonChan)

	defer handler.Close()

	go func() {
		fmt.Println("Listening for Dash buttons...")
		if err := handler.Listen(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case button := <-buttonChan:
			log.Printf("Received registered button press: %s", button)
			return // Exit main and close handler
		case button := <-unknownButtonChan:
			log.Printf("Unknown dash button pressed. MAC address: %s\n", button)
		}
	}
}
