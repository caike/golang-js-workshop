package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

var serverURL = flag.String("server", "http://localhost:3000", "Web API")
var deviceName = flag.String("deviceName", "dashboard-device", "Name of the device")

var endPointForPOST bytes.Buffer
var interrupt = make(chan os.Signal, 1)

func main() {
	daemon.SdNotify(false, "READY=1")
	// Terminates the program upon `ctrl+c`
	signal.Notify(interrupt, os.Interrupt)

	flag.Parse()
	sendData()
}

func sendData() {
	done := make(chan struct{})
	output := make(chan string)

	url := *serverURL
	endPointForPOST.WriteString(url)
	endPointForPOST.WriteString("/data")

	go func() {
		defer close(done)

		for {
			// Hardcoded command for now.
			// TODO: read the command from a config file.
			commandStr := "ls ~/"
			out, err := exec.Command("bash", "-c", commandStr).Output()
			if err != nil {
				log.Fatal(err)
			}
			output <- string(string(out))
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		select {
		case commandOutput := <-output:
			sendStatusToServer(commandOutput)
			daemon.SdNotify(false, "WATCHDOG=1")
		case <-interrupt:
			log.Println("Closing the boxtuneProcessGuard client")
			os.Exit(0)
			return
		}
	}
}

func sendStatusToServer(output string) {
	status := deviceStatus{DeviceName: *deviceName, CommandOutput: output}
	jsonToSend, jsonErr := json.Marshal(status)
	if jsonErr != nil {
		log.Fatal("Error marshalling JSON")
	}
	jsonStr := []byte(jsonToSend)
	req, err := http.NewRequest("POST", endPointForPOST.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal("Erro: ", err)
	}
	req.Close = true
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Device-Name", *deviceName)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error posting status. Error: ", err)
		// must return in order to skip resp.Body.Close()
		// which causes dereference error
		return
	}
	defer resp.Body.Close()
	if resp.Status != "201 Created" {
		log.Println("Erro POSTing duplicates", resp.Status)
	}
	log.Println("Posted")
}

type deviceStatus struct {
	DeviceName    string
	CommandOutput string
}
