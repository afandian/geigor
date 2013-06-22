package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"
)

// HTML buffers
var monitorHtmlBuffer *bytes.Buffer = &bytes.Buffer{}
var indexBuffer *bytes.Buffer
var demoBuffer *bytes.Buffer
var collectorJsBuffer *bytes.Buffer = &bytes.Buffer{}

// Config
var portNumber string
var servedFrom string

// Template args
type HtmlConfig struct {
	Port       string
	ServedFrom string
}

// Concurrent bits.

// The count for the current buffer.
var currentCount int = 0

// Duration of the buffer
var bufferDuration = 1 * time.Second

// A channel to regulate change of the current buffer's count.
var tickCommands = make(chan int, 100)

// A channel to register or de-register websocket listener sync messages.
var websocketSyncRegister = make(chan registrationInstruction, 100)

// A collection of channels which are used as sync signals to websocket handlers.
// The channel is sent the current count at the time of sending.
var syncSignals map[chan int]bool = make(map[chan int]bool)

// These values are sent down the tick channel.
const (
	// Increment buffer.
	TICK = 0

	// Reset buffer.
	RESET = 1
)

// Instruction sent in registration.
const (
	REGISTER    = 0
	DE_REGISTER = 1
)

// A register / deregister message.
// Sent with a channel which is used to signal to the monitor server websocket.
type registrationInstruction struct {
	instruction int
	signal      chan int
}

// Keep the current buffer up to date by recieving ticks from web socket handlers and the timer.
func ticker() {
	var command int

	for {
		command = <-tickCommands
		if command == TICK {
			currentCount++
		} else if command == RESET {
			currentCount = 0
		}
	}
}

// Send periodic resets each buffer length.
func timer() {
	for {
		time.Sleep(bufferDuration)

		// Signal the current count to all websockets.
		for signal, _ := range syncSignals {
			signal <- currentCount
		}

		// Then reset the counter.
		tickCommands <- RESET
	}
}

// Register and de-register web socket threads who listen out for sync ticks.
// Websockets register and de-register themselves into syncSignals.
func registrar() {
	var instruction registrationInstruction
	for {
		instruction = <-websocketSyncRegister

		if instruction.instruction == REGISTER {
			syncSignals[instruction.signal] = true
		} else if instruction.instruction == DE_REGISTER {
			delete(syncSignals, instruction.signal)
		}
	}
}

// Echo the data received on the WebSocket.
func tickServer(ws *websocket.Conn) {
	defer ws.Close()

	// A channel for this websocket connection.
	// It will recieve ticks with the current count, to broadcast back.
	var myChannel = make(chan int)

	// Register me.
	websocketSyncRegister <- registrationInstruction{REGISTER, myChannel}

	var count int

	// Every tick we'll get a message down the channel with the current count.
	for true {
		count = <-myChannel
		_, err := ws.Write([]byte(fmt.Sprintf("%d", count)))

		// If the socket was closed, deregister.
		if err != nil {
			websocketSyncRegister <- registrationInstruction{DE_REGISTER, myChannel}
			return
		}
	}
}

func monitor(response http.ResponseWriter, r *http.Request) {
	_, err := response.Write(monitorHtmlBuffer.Bytes())

	if err != nil {
		panic(err)
	}
}

func index(response http.ResponseWriter, r *http.Request) {
	_, err := response.Write(indexBuffer.Bytes())

	if err != nil {
		panic(err)
	}

}

func demo(response http.ResponseWriter, r *http.Request) {
	_, err := response.Write(demoBuffer.Bytes())

	if err != nil {
		panic(err)
	}

}

func collectorJs(response http.ResponseWriter, r *http.Request) {
	_, err := response.Write(collectorJsBuffer.Bytes())

	if err != nil {
		panic(err)
	}

}

func collectorEndpoint(response http.ResponseWriter, r *http.Request) {
	tickCommands <- TICK
}

func printHelp() {
	fmt.Println("Arguments: <served from> <port number> <option>*")
	fmt.Println("Options: index, demo")
}

// This example demonstrates a trivial echo server.
func main() {

	// Get args.
	if len(os.Args) < 3 {
		printHelp()
		os.Exit(1)
	}

	// The full URL where this is served.
	servedFrom = os.Args[1]

	// Port number is a string as this is passed to the server.
	// ListenAndServe will throw an error if this is bad.
	portNumber := os.Args[2]

	// Start the state tracking bits.
	go ticker()
	go timer()
	go registrar()

	for _, arg := range os.Args {
		if arg == "index" {

			// The the index page.
			http.HandleFunc("/", index)
		}
	}

	for _, arg := range os.Args {
		if arg == "demo" {

			// The the demo page.
			http.HandleFunc("/demo", demo)
		}
	}

	// Render templates.
	templateConfig := HtmlConfig{portNumber, servedFrom}

	monitorHtmlTemplate, err := template.New("a").Parse(MONITOR_HTML)
	if err != nil {
		panic(err)
	}

	collectorJsTemplate, err := template.New("b").Parse(COLLECTOR_JS)
	if err != nil {
		panic(err)
	}

	monitorHtmlTemplate.Execute(monitorHtmlBuffer, templateConfig)
	indexBuffer = bytes.NewBufferString(INDEX_HTML)
	demoBuffer = bytes.NewBufferString(DEMO_HTML)
	collectorJsTemplate.Execute(collectorJsBuffer, templateConfig)

	// The endpoint that provides a monitoring websocket.
	http.Handle("/monitor-endpoint", websocket.Handler(tickServer))

	// The monitor HTML interface.
	http.HandleFunc("/monitor", monitor)

	// The collection JavaScript to put in the target page.
	http.HandleFunc("/collector.js", collectorJs)

	// The endpoint that recieves ticks.
	http.HandleFunc("/collector-endpoint", collectorEndpoint)

	err = http.ListenAndServe("0.0.0.0:"+portNumber, nil)
	if err != nil {

		fmt.Println("Error starting server: " + err.Error())
		printHelp()
		os.Exit(1)
	}
}
