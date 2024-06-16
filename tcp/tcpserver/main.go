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
		fmt.Println("Please provide host:port")
		os.Exit(1)
	}

	// Resolve the string address to a TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start listening for TCP connections on the given address
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	start := time.Now() // Start measuring time after setting up listener

	fmt.Printf("TCP server started and listening on %s\n", tcpAddr.String())

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Handle new connections in a Goroutine for concurrency
		go handleConnection(conn)
	}

	elapsed := time.Since(start) // Calculate total elapsed time
	fmt.Printf("Total execution time: %v\n", elapsed)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Read from the connection until a new line is received
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the data read from the connection to the terminal
		fmt.Print("> ", string(data))

		// Write back the same message to the client
		conn.Write([]byte("Hello TCP Client\n"))
	}
}
