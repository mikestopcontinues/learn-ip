package main

import "sync"

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
	3389: "RDP",
	5900: "VNC",
}

func main() {
	var wg sync.WaitGroup
	subnet := "192.168.1"

	// step 1 here

	wg.Wait()
}

func scanCommonPorts(ip string) {
	// step 2 here
}
