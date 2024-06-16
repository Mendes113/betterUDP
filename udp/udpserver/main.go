package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide host:port")
		os.Exit(1)
	}

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("UDP server started and listening on %s\n", udpAddr.String())

	// Read from UDP listener in endless loop
	for {
		var buf [512]byte
		// Start measuring time for receiving message
		recvStart := time.Now()

		_, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		recvElapsed := time.Since(recvStart) // Calculate time taken for receiving

		fmt.Print("> ", string(buf[0:]))

		// Start measuring time for sending response
		sendStart := time.Now()

		// Write back the message over UDP
		_, err = conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		sendElapsed := time.Since(sendStart) // Calculate time taken for sending

		// Print time spent for receiving and sending
		fmt.Printf("Time spent receiving: %v\n", recvElapsed)
		fmt.Printf("Time spent sending: %v\n", sendElapsed)
	}
}
