package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	WINDOW_SIZE = 10
)

var (
	nextSeqNum   int
	baseSeqNum   int
	congestion   bool
	receiveQueue []string
	mutex        sync.Mutex
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORT := ":" + arguments[1]

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		msg := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received: %s\n", msg)

		if msg == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		mutex.Lock()
		receiveQueue = append(receiveQueue, msg)
		mutex.Unlock()

		// Simulating congestion
		if congestion {
			time.Sleep(time.Millisecond * 100)
			congestion = false
		}

		// Process message and send ACK
		if nextSeqNum < baseSeqNum+WINDOW_SIZE && nextSeqNum < len(receiveQueue) {
			expectedSeqNum := strconv.Itoa(nextSeqNum + 1)
			ack := []byte(expectedSeqNum)
			_, err = connection.WriteToUDP(ack, addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			nextSeqNum++
		} else {
			congestion = true
			fmt.Println("Congestion detected! Waiting...")
		}

		// Update base sequence number
		mutex.Lock()
		if len(receiveQueue) > 0 && nextSeqNum >= baseSeqNum+len(receiveQueue) {
			baseSeqNum++
			receiveQueue = receiveQueue[1:]
		}
		mutex.Unlock()
	}
}
