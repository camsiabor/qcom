package qnet

import "net"

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
