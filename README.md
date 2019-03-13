# go-dash-button
[![GoDoc](https://godoc.org/github.com/ahornerr/dash-button?status.svg)](https://godoc.org/github.com/ahornerr/dash-button)

This Go library listens for DHCP requests and calls registered handlers when a request is made.
Typically this will be used for Amazon Dash buttons which, when pressed, connect to WiFI and send a DHCP request.
Devices are identified by their hardware address (MAC address) and handlers are registered as so.

## Example
```golang
// Make sure to receive from your button handler channels. A non blocking channel send is used in the library.
// buttonChan := make(dashbutton.ButtonHandler)
// unknownButtonChan := make(dashbutton.ButtonHandler)

func listen(buttonChan, unknownButtonChan dashbutton.ButtonHandler) {
    handler := dashbutton.NewHandler()
    handler.SetUnknownButtonHandler(unknownButtonChan)
    handler.AddButtonHandler("fc:65:de:b2:8c:df", buttonChan)

    defer handler.Close()

    fmt.Println("Listening for Dash buttons...")
    
    // handler.Listen() blocks until an error is returned or handler.Close() is called
    err := handler.Listen();
}
```

## Docker
Given the `-p 0.0.0.0:67:67/udp` flag to `docker run`, this library (and example app) is usable inside a Docker container with no modifications or host network flag necessary.

```bash
$ go build -o dash-button-example example/main.go

$ docker build -t dash-button-example -f example.Dockerfile .
...
Successfully built 292fa312f9bc
Successfully tagged dash-button-example:latest

$ docker run -p 0.0.0.0:67:67/udp dash-button-example
Listening for Dash buttons...
2019/03/13 20:36:22 Received registered button press: fc:65:de:b2:8c:df
```