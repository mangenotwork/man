package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	getMac()

	// getIP()
}

func getMac() {
	iFaces, _ := net.Interfaces()
	for _, v := range iFaces {
		//log.Println(v.Name, v.HardwareAddr.String())

		addrs, err := v.Addrs()
		if err != nil {
			fmt.Println("Failed to get addresses for interface", v.Name, ":", err)
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				fmt.Println("Interface:", v.Name)
				fmt.Println("IP Address:", ipNet.IP)
				fmt.Println("MAC Address:", v.HardwareAddr.String())
				fmt.Println("______________________")
			}
		}

	}
}

func getIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local ip failed")
	}
	for _, address := range addrs {
		//log.Println(address)
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			log.Println(ipnet)
			if ipnet.IP.To4() != nil {
				log.Println(ipnet.IP.String())
			}
		}
	}
}
