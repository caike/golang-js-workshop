package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"

	"strings"

	"github.com/gorilla/websocket"
)

// Takes destination server as single argument,
// otherwise uses localhost:8080 as default
var addr = flag.String("server", "localhost:3000", "Websocket Server")

func main() {

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	connectHeader := make(http.Header)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/pi"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), connectHeader)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	output := make(chan string)

	go func() {
		defer c.Close()
		defer close(done)

		for {
			_, commandFromWeb, commandErr := c.ReadMessage()
			if commandErr != nil {
				log.Fatal("Error reading command")
			}
			output <- string(commandFromWeb)
		}
	}()

	for {
		select {
		case command := <-output:

			commandSlice := strings.Split(command, "::")

			out, err := exec.Command("bash", "-c", commandSlice[1]).Output()
			if err != nil {
				log.Fatal(err)
			}

			var responseBuffer bytes.Buffer
			responseBuffer.WriteString(commandSlice[0]) // clientName
			responseBuffer.WriteString("::")            // separator
			responseBuffer.WriteString(string(out))     // command output

			writeErr := c.WriteMessage(websocket.TextMessage, []byte(responseBuffer.String()))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
		case <-interrupt:
			log.Println("interrupt signal, closing the client")
			c.Close()
			return
		}
	}
}
