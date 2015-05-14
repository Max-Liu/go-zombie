package zombie

import (
	"errors"
	"net"
	"os"
)

func GetIp() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		ip := net.ParseIP(a)
		if ip.IsGlobalUnicast() {
			return ip.String(), nil
		}
	}
	return "", errors.New("do not have ip address")
}
