package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/stianeikeland/go-rpio"
	"github.com/yryz/ds18b20"
)

var pin rpio.Pin
var pinNumber int
var testmode bool

func main() {
	testmode = false
	data, _ := ioutil.ReadFile("./config")
	s := strings.Trim(string(data), "\n")
	pinNumber, _ := strconv.Atoi(s)
	pin = rpio.Pin(pinNumber)

	err := rpio.Open()
	if err != nil {
		fmt.Println("unable to open gpio", err.Error())
		fmt.Println("running in test mode")
		testmode = true
	} else {
		fmt.Println("creating channel")
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			cleanup(pin, pinNumber)
			os.Exit(1)
		}()

		sensors, err := ds18b20.Sensors()
		if err != nil {
			panic(err)
		}

		fmt.Printf("sensor IDs: %v\n", sensors)

		for _, sensor := range sensors {
			t, err := ds18b20.Temperature(sensor)
			if err == nil {
				fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
			}
		}

	}
	fmt.Println("Setting up http handlers")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/config", configHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func configHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}

func cleanup(pin rpio.Pin, pinNumber int) {
	fmt.Println("Cleaning up pin", pinNumber)
	pin.Write(rpio.Low)
	rpio.Close()
}
