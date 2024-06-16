package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	WINDOW_SIZE = 1
)

var (
	nextSeqNum            int
	baseSeqNum            int
	congestion            bool
	receiveQueue          []string
	mutex                 sync.Mutex
	packetRate            float64
	lastPacketTime        time.Time
	CONGESTION_THRESHOLD  = 1000                         // Initial congestion queue limit
	CONGESTION_DELAY      = 1000 * time.Millisecond     // Initial congestion delay time
	CSV_FILE_PATH         = "./times.csv"
	startTime, endTime    time.Time
)

func Server(PORT string) {
	startTime = time.Now()
	

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
	lastPacketTime = time.Now()

	go AdjustCongestionParameters()

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		go HandlePacket(connection, buffer[:n], addr)
	}
}

func HandlePacket(connection *net.UDPConn, packet []byte, addr *net.UDPAddr) {


	

	

	msg := strings.TrimSpace(string(packet))

	color.Magenta("Received: %s\n", msg)
	if msg == "STOP" {
		fmt.Println("Exiting UDP server!")
		
		return
	}

	// Simulating congestion
	if congestion {
		time.Sleep(CONGESTION_DELAY) // Use the dynamically adjusted congestion delay
		congestion = false
	}

	mutex.Lock()
	defer mutex.Unlock()

	seqNum, err := strconv.Atoi(strings.Split(msg, " ")[1])
	if err != nil {
		fmt.Println("Invalid message format")
		return
	}

	if seqNum == baseSeqNum+1 {
		receiveQueue = append(receiveQueue, msg)
		baseSeqNum++
		expectedSeqNum := strconv.Itoa(baseSeqNum)
		ack := []byte(expectedSeqNum)
		_, err = connection.WriteToUDP(ack, addr)
		if err != nil {
			fmt.Println(err)
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
}

func AdjustCongestionParameters() {
	for {
		time.Sleep(1 * time.Second) // Adjust parameters every 5 seconds
		mutex.Lock()

		// Adjust congestion threshold based on packet arrival rate
		if packetRate > 20 {
			CONGESTION_THRESHOLD = 100
			CONGESTION_DELAY = 200 * time.Millisecond
		} else if packetRate > 10 {
			CONGESTION_THRESHOLD = 75
			CONGESTION_DELAY = 400 * time.Millisecond
		} else if packetRate > 5 {
			CONGESTION_THRESHOLD = 50
			CONGESTION_DELAY = 600 * time.Millisecond
		} else {
			CONGESTION_THRESHOLD = 25
			CONGESTION_DELAY = 1000 * time.Millisecond
		}

		color.Magenta("Adjusted CONGESTION_THRESHOLD to %d and CONGESTION_DELAY to %v based on packet rate: %.2f packets/sec\n", CONGESTION_THRESHOLD, CONGESTION_DELAY, packetRate)
		mutex.Unlock()
	}
}


