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

	// Resolve the string address to a TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connect to the address with TCP
	start := time.Now() // Start measuring time after resolving address
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send a message to the server
	_, err = conn.Write([]byte("Hello TCP Server\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read from the connection until a new line is received
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data read from the connection to the terminal
	fmt.Print("> ", string(data))

	elapsed := time.Since(start) // Calculate total elapsed time
	fmt.Printf("Total execution time: %v\n", elapsed)
}
