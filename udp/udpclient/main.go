package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide host:port to connect to")
		os.Exit(1)
	}

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// Start measuring time for sending messages
	sendStart := time.Now()
	for i := 1; i <= 100; i++ {
		// Send a message to the server
		_, err = conn.Write([]byte("Hello UDP Server\n"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	sendElapsed := time.Since(sendStart) // Calculate time taken for sending

	fmt.Println("Messages sent to UDP server")

	// Start measuring time for receiving response
	recvStart := time.Now()

	// Read from the connection until a new line is received
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	recvElapsed := time.Since(recvStart) // Calculate time taken for receiving

	// Print the data read from the connection to the terminal
	fmt.Print("> ", string(data))

	// Calculate total time spent
	totalElapsed := sendElapsed + recvElapsed

	// Print time spent for sending, receiving, and total time
	fmt.Printf("Time spent sending: %v\n", sendElapsed)
	fmt.Printf("Time spent receiving: %v\n", recvElapsed)
	fmt.Printf("Total time spent: %v\n", totalElapsed)
}
