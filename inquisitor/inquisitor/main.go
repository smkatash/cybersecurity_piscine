package main

import (
	"log"
	"os"
	"net"
)

type DeviceId struct {
	ip net.IP
	mac net.HardwareAddr
}

type InputDevices struct {
	source DeviceId
	target DeviceId
}

func main() {
	if (len(os.Args) != 5) {
		log.Fatal("usage: ip-source mac-source ip-target mac-target")
	}
	input := InputDevices{}
	ParseAddresses(&input) 
	go Monitor(&input)
	ArpSpoof(&input)
}
