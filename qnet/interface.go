package qnet

import (
	"fmt"
	"net"
	"strings"
)

func AllNetInterfaceIP() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// handle err
	var ips []net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ips, err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				if ips == nil {
					ips = []net.IP{ip}
				} else {
					ips = append(ips, ip)
				}
			}
		}
	}
	return ips, nil
}

func AllNetInterfaceIPString() ([]string, error) {
	var ips, err = AllNetInterfaceIP()
	var ipstrs []string
	if ips != nil {
		ipstrs = make([]string, len(ips))
		for i, ip := range ips {
			ipstrs[i] = ip.String()
		}
	}
	return ipstrs, err
}

func CurrentIPAndMac() (string, string, error) {

	//----------------------
	// Get the local machine IP address
	// https://www.socketloop.com/tutorials/golang-how-do-I-get-the-local-ip-non-loopback-address
	//----------------------

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", "", err
	}

	var currentIP, currentNetworkHardwareName string

	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}

	// get all the system's or local machine's network interfaces
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				// only interested in the name with current IP address
				if strings.Contains(addr.String(), currentIP) {
					currentNetworkHardwareName = interf.Name
				}
			}
		}
	}

	// extract the hardware information base on the interface name
	// capture above
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	fmt.Println("Hardware name : ", name)
	fmt.Println("MAC address : ", macAddress)

	// verify if the MAC address can be parsed properly
	_, err = net.ParseMAC(macAddress.String())

	if err != nil {
		return "", "", err
	}

	return currentIP, macAddress.String(), nil

}
