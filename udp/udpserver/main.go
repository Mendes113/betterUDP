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

	messageCounter := 0
	totalTime := time.Duration(0)

	// Read from UDP listener in loop until we receive 1000 messages
	for messageCounter < 1000 {
		var buf [2048]byte
		// Start measuring time for receiving message
		recvStart := time.Now()

		_, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		recvElapsed := time.Since(recvStart) // Calculate time taken for receiving

		fmt.Print("> ", string(buf[0:]))

		// // Start measuring time for sending response
		// sendStart := time.Now()

		// Write back the message over UDP
		_, err = conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// sendElapsed := time.Since(sendStart) // Calculate time taken for sending

		// // // Print time spent for receiving and sending
		// fmt.Printf("Time spent receiving: %v\n", recvElapsed)
		// fmt.Printf("Time spent sending: %v\n", sendElapsed)

		// Update counters and time
		messageCounter++
		// log.Printf("counter %d ", messageCounter)
		totalTime += recvElapsed
	}

	// Print total time spent receiving 1000 messages
	fmt.Printf("Total time spent receiving 1000 messages: %v\n", totalTime)
}
