package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mikehamer/crazyserver/cache"
	"github.com/mikehamer/crazyserver/crazyflie"
	"github.com/mikehamer/crazyserver/crazyradio"
)

func main() {
	var err error
	channel := flag.Uint("channel", 80, "Radio channel")
	address := flag.Uint64("address", 0xE7E7E7E701, "Radio address")
	flag.Parse()
	cache.Init()

	radio, err := crazyradio.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer radio.Close()

	radio.SetChannel(uint8(*channel))

	cf, err := crazyflie.Connect(radio, *address)
	if err != nil {
		log.Fatal(err)
	}
	defer cf.Disconnect()
	// log.Println("Rebooting")
	// cf.RebootToFirmware()
	// log.Println("Rebooted")

	// <-time.After(1 * time.Second)

	cf.LogSystemReset()
	err = cf.LogTOCGetList()
	if err != nil {
		log.Fatal(err)
	}

	err = cf.ParamTOCGetList()
	if err != nil {
		log.Fatal(err)
	}

	val, err := cf.ParamRead("kalman.pNAcc_xy")
	fmt.Println(val)
	err = cf.ParamWrite("kalman.pNAcc_xy", float32(3.14159))
	val, err = cf.ParamRead("kalman.pNAcc_xy")
	fmt.Println(val)

	// Unlock commander
	cf.SetpointSend(0, 0, 0, 0)
	// Commander packets needs to be sent at regular interval, otherwise the
	// commander watchdog will cut the motors
	stop := false
	go func() {
		for !stop {
			cf.SetpointSend(0, 0, 0, 4000)
			<-time.After(20 * time.Millisecond)
		}
	}()
	<-time.After(5 * time.Second)
	stop = true
	<-time.After(40 * time.Millisecond)
	cf.SetpointSend(0, 0, 0, 0)
	<-time.After(1 * time.Second)
}
