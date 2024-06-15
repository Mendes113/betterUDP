package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	WINDOW_SIZE = 1
)

var (
	nextSeqNum        int
	baseSeqNum        int
	congestion        bool
	receiveQueue      []string
	mutex             sync.Mutex
	packetRate        float64
	lastPacketTime    time.Time
	CONGESTION_THRESHOLD = 750        // Inicialmente, o limite de mensagens na fila para ativar congestionamento
	CONGESTION_DELAY      = 1000 * time.Millisecond // Inicialmente, o tempo de atraso para simular congestionamento
)

func Server(PORT string) {
	

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
	now := time.Now()
	elapsed := now.Sub(lastPacketTime).Seconds()
	lastPacketTime = now

	// Atualiza a taxa de chegada de pacotes (simples média móvel)
	packetRate = (packetRate + (1 / elapsed)) / 2

	msg := strings.TrimSpace(string(packet))
	fmt.Printf("Received: %s\n", msg)

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
		time.Sleep(1 * time.Second) // Ajusta os parâmetros a cada 5 segundos
		mutex.Lock()

		// Ajusta o limiar de congestionamento com base na taxa de chegada de pacotes
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

		fmt.Printf("Adjusted CONGESTION_THRESHOLD to %d and CONGESTION_DELAY to %v based on packet rate: %.2f packets/sec\n", CONGESTION_THRESHOLD, CONGESTION_DELAY, packetRate)
		mutex.Unlock()
	}
}
