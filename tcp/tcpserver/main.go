package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var mutex sync.Mutex
var CSV_FILE_PATH = "/home/mendes/Documents/Github/betterUDP/execution_times.csv"
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

	// Channel to signal server shutdown after 1000 messages
	stopCh := make(chan struct{})

	// Handle OS signals in a separate Goroutine
	go func() {
		<-signalCh
		fmt.Println("\nReceived interrupt signal. Shutting down gracefully...")
		listener.Close()
		printExecutionTime(start)
		os.Exit(0)
	}()

	// Accept connections and handle them
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-stopCh:
					// Server shutdown signal received
					return
				default:
					fmt.Println(err)
					continue
				}
			}

			// Handle new connections in a Goroutine for concurrency
			go handleConnection(conn, start, stopCh)
		}
	}()

	// Block until shutdown signal is received
	<-stopCh
	printExecutionTime(start)
}

func handleConnection(conn net.Conn, start time.Time, stopCh chan struct{}) {
	defer conn.Close()

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Print execution time when EOF is encountered
				printExecutionTime(start)
			} else {
				fmt.Println(err)
			}
			return
		}

		// Increment message count
		mutex.Lock()
		messageCount++
		mutex.Unlock()

		fmt.Print("> ", string(data))

		conn.Write([]byte("Hello TCP Client\n"))

		if messageCount >= MAX_MESSAGES {
			// Signal the server to stop accepting new connections
			close(stopCh)
			return
		}
	}
}

func printExecutionTime(start time.Time) {
	elapsed := time.Since(start).Milliseconds()
	fmt.Printf("Total execution time: %d ms\n", elapsed)
	// err := LogElapsedTime(elapsed)
	// if err != nil {
	// 	fmt.Println("Error logging elapsed time:", err)
	// }
}

// func LogElapsedTime(serverTimeMs int64) error {
// 	// Open CSV file in read mode
// 	file, err := os.OpenFile(CSV_FILE_PATH, os.O_RDWR, 0644)
// 	if err != nil {
// 		return fmt.Errorf("error opening CSV file: %w", err)
// 	}
// 	defer file.Close()

// 	// Read existing CSV records
// 	reader := csv.NewReader(file)
// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		return fmt.Errorf("error reading CSV file: %w", err)
// 	}

// 	// Open CSV file again in write mode to overwrite with updated data
// 	file, err = os.OpenFile(CSV_FILE_PATH, os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return fmt.Errorf("error opening CSV file for writing: %w", err)
// 	}
// 	defer file.Close()

// 	// Create CSV writer
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	// Update existing records with new time in the third column
// 	for _, record := range records {
// 		if len(record) >= 2 {
// 			// Append serverTimeMs to the end of each record
// 			record = append(record, strconv.FormatInt(serverTimeMs, 10))
// 			if err := writer.Write(record); err != nil {
// 				return fmt.Errorf("error writing updated record to CSV: %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }
