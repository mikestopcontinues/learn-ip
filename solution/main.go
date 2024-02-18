package main

import (
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
)

var commonPorts = map[int]string{
	20:   "FTP Data Transfer",
	21:   "FTP Command Control",
	22:   "SSH",
	23:   "Telnet",
	25:   "SMTP",
	53:   "DNS",
	80:   "HTTP",
	110:  "POP3",
	143:  "IMAP",
	443:  "HTTPS",
	587:  "SMTP (Email submission)",
	3306: "MySQL",
	3389: "RDP",
	5900: "VNC",
}

func scanCommonPorts(ip string) {
	var wg sync.WaitGroup

	for port, description := range commonPorts {
		wg.Add(1)

		go func(port int, description string) {
			defer wg.Done()

			address := fmt.Sprintf("%s:%d", ip, port)

			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err != nil {
				return
			}
			conn.Close()

			fmt.Printf("%s:%d (%s)\n", ip, port, description)
		}(port, description)
	}

	wg.Wait()
}

func main() {
	subnet := "192.168.1"
	var wg sync.WaitGroup

	for i := 1; i <= 254; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			ip := fmt.Sprintf("%s.%d", subnet, i)

			_, err := exec.Command("ping", "-c", "1", "-W", "1", ip).CombinedOutput()
			if err != nil {
				return
			}

			fmt.Println(ip)

			scanCommonPorts(ip)
		}(i)
	}

	wg.Wait()
}
