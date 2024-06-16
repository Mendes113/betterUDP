package main

import (
	"bufio"
	
	"fmt"
	"net"
	"os"
	
	"sync"
	"time"
)

var mutex sync.Mutex 
var CSV_FILE_PATH = "/home/mendes/Documents/Github/betterUDP/execution_times.csv"

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

	// Loop to send 1000 messages
	for i := 1; i <= 1000; i++ {
		message := fmt.Sprintf("Message %d\n", i)

		// Send a message to the server
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Read response from the server
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the data read from the connection to the terminal
		fmt.Print("> ", string(data))
	}

	elapsed := time.Since(start) // Calculate total elapsed time
	fmt.Printf("Total execution time: %v\n", elapsed)
	// LogElapsedTime(elapsed)
}

// func LogElapsedTime(serverTime time.Duration) error {
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

// 	// Open CSV file again in append mode to add new records
// 	file, err = os.OpenFile(CSV_FILE_PATH, os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return fmt.Errorf("error opening CSV file for writing: %w", err)
// 	}
// 	defer file.Close()

// 	// Create CSV writer
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	// Determine if we need to start a new line
// 	startNewLine := true

// 	// Iterate over existing records to check if we need to start a new line
// 	for _, record := range records {
// 		if len(record) >= 3 {
// 			// Check if the last column is empty or needs to be skipped
// 			if record[2] == "" {
// 				// Update the third column with serverTimeMs
// 				record[2] = strconv.FormatInt(serverTime.Milliseconds(), 10)
// 				if err := writer.Write(record); err != nil {
// 					return fmt.Errorf("error writing updated record to CSV: %w", err)
// 				}
// 				startNewLine = false // No need to start a new line
// 				break                // Exit loop after updating
// 			}
// 		}
// 	}

// 	// If no suitable place was found, start a new line
// 	if startNewLine {
// 		newRecord := []string{"", "", strconv.FormatInt(serverTime.Milliseconds(), 10)}
// 		if err := writer.Write(newRecord); err != nil {
// 			return fmt.Errorf("error writing new record to CSV: %w", err)
// 		}
// 	}

// 	return nil
// }
