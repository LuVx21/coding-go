package web

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func WebProvider() map[string]any {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	ip, err := GetOutBoundIP()
	if err != nil {
		ip = ""
	}
	return map[string]any{
		"host": hostname,
		"ip":   ip,
	}
}

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
