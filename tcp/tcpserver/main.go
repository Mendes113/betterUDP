package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var mutex sync.Mutex
var CSV_FILE_PATH = "execution_times.csv"
var messageCount = 0
const MAX_MESSAGES = 1000

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

	// Channel to receive OS signals (e.g., SIGINT, SIGTERM)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Handle OS signals in a separate Goroutine
	go func() {
		<-signalCh
		fmt.Println("\nReceived interrupt signal. Shutting down gracefully...")
		listener.Close()
	}()

	// Accept connections and handle them
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Increment message count
		messageCount++

		// Handle new connections in a Goroutine for concurrency
		go func() {
			defer conn.Close()

			for {
				data, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Print("> ", string(data))

				conn.Write([]byte("Hello TCP Client\n"))

				if messageCount >= MAX_MESSAGES {
					// Calculate total elapsed time
					elapsed := time.Since(start)
					fmt.Printf("Received all 1000 messages. Total execution time: %v\n", elapsed)

					// Log the elapsed time in CSV
					err := LogElapsedTime(elapsed.Milliseconds(), 0) // Client time is not relevant here
					if err != nil {
						fmt.Println("Error logging elapsed time:", err)
					}

					// Close listener to stop accepting new connections
					listener.Close()
					return
				}
			}
		}()
	}
}

func LogElapsedTime(serverTimeMs, clientTimeMs int64) error {
	// Open CSV file in append mode
	file, err := os.OpenFile(CSV_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write execution times to CSV
	mutex.Lock()
	defer mutex.Unlock()

	if err := writer.Write([]string{strconv.FormatInt(serverTimeMs, 10), strconv.FormatInt(clientTimeMs, 10)}); err != nil {
		return fmt.Errorf("error writing execution times to CSV: %w", err)
	}

	return nil
}
