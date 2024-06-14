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
	WINDOW_SIZE             = 1
	CONGESTION_DELAY        = 1000 * time.Millisecond // Tempo de atraso aumentado para simular congestionamento
	CONGESTION_THRESHOLD    = 2       // Limite de mensagens na fila para ativar congestionamento
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

		// Simulating congestion
		if congestion {
			time.Sleep(CONGESTION_DELAY) // Use the increased congestion delay
			congestion = false
		}

		mutex.Lock()
		seqNum, err := strconv.Atoi(strings.Split(msg, " ")[1])
		if err != nil {
			fmt.Println("Invalid message format")
			mutex.Unlock()
			continue
		}

		if seqNum == baseSeqNum+1 {
			receiveQueue = append(receiveQueue, msg)
			baseSeqNum++
			expectedSeqNum := strconv.Itoa(baseSeqNum)
			ack := []byte(expectedSeqNum)
			_, err = connection.WriteToUDP(ack, addr)
			if err != nil {
				fmt.Println(err)
				mutex.Unlock()
				return
			}
			fmt.Printf("Sent ACK: %s\n", expectedSeqNum)
		} else if seqNum > baseSeqNum {
			fmt.Println("Message out of order, dropping:", msg)
		}

		// Simulate random congestion
		if len(receiveQueue) > CONGESTION_THRESHOLD {
			fmt.Println("Congestion detected! Slowing down...")
			congestion = true
		}

		mutex.Unlock()
	}
}
