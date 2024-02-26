package main

import (
	"net"
	"os"
	"log"
)


func ParseIPAddress(ipAddress string) []byte {
	ip := net.ParseIP(ipAddress)
	if ip != nil {
		return net.IP(ip).To4() 
	}
	return nil
}

func ParseMACAddress(macAddress string) []byte {
	hardWareAddr, err := net.ParseMAC(macAddress)
	if err != nil {
		return nil
	}
	return hardWareAddr
}


func ParseAddresses(input *InputDevices) {
	args := os.Args
	
	sourceDeviceId := DeviceId{
		ip: ParseIPAddress(args[1]),              
		mac: ParseMACAddress(args[2]),
	}
	targetDeviceId := DeviceId{
		ip: ParseIPAddress(args[3]),              
		mac: ParseMACAddress(args[4]),
	}

	if sourceDeviceId.ip == nil && targetDeviceId.ip == nil {
		log.Fatal("only ip4 is supported")
	}

	if sourceDeviceId.mac == nil || targetDeviceId.mac == nil {
		log.Fatal("invalid mac address")
	}

	input.source = sourceDeviceId
	input.target = targetDeviceId
}