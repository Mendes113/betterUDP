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


var (
	// Variables to handle congestion control
	nextSeqNum            int
	baseSeqNum            int
	congestion            bool
	receiveQueue          []string
	mutex                 sync.Mutex
	packetRate            float64
	lastPacketTime        time.Time
	CONGESTION_THRESHOLD  = 1000                         // Congestao threshold incial 
	CONGESTION_DELAY      = 1000 * time.Millisecond     // atrazo de congestionamento inicial
	CSV_FILE_PATH         = "./times.csv"
	startTime, endTime    time.Time
	WINDOW_SIZE           = 2 // Tamanho da janela de congestionamento inicial
)

func Server(PORT string) {
	startTime = time.Now()
	
	// Resolve UDP address
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Listen on UDP port
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Close connection when application ends
	defer connection.Close()
	// Buffer para armazenar os dados recebidos
	buffer := make([]byte, 1024)
	lastPacketTime = time.Now()
	// Iniciar goroutine para ajustar os parâmetros de congestionamento
	// go routine é similar a uma  thread
	// entretanto em go uma go routine é mais leve que uma thread
	// diversas go routines podem ser executadas em uma única thread
	go AdjustCongestionParameters()

	for {
		// Read data from UDP connection
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Calcula a taxa de pacotes recebidos
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

	// Simula congestionamento
	if congestion {
		time.Sleep(CONGESTION_DELAY) // Usa o tempo de atraso de congestionamento
		congestion = false
	}
	//mutex serve para garantir que apenas uma go routine execute o código por vez
	mutex.Lock()
	//defer garante que a função seja chamada no final da execução da função atual
	defer mutex.Unlock()
	// se o tempo decorrido desde o último pacote for maior que 1 segundo
	seqNum, err := strconv.Atoi(strings.Split(msg, " ")[1])
	if err != nil {
		fmt.Println("Invalid message format")
		return
	}
	// se o número de sequência do pacote for maior que o próximo número de sequência esperado
	if seqNum == baseSeqNum+1 {
		// Adiciona o pacote à fila de recebimento
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

	// Simula congestionamento
	if len(receiveQueue) > CONGESTION_THRESHOLD {
		fmt.Println("Congestion detected! Slowing down...")
		congestion = true
	}
}

func AdjustCongestionParameters() {
	// Ajusta os parâmetros de congestionamento
	for {
		// ajusta a cada segundo
		time.Sleep(1 * time.Second) 
		//mutex serve para garantir que apenas uma go routine execute o código por vez
		mutex.Lock()
		
		// ajusta a janela de congestionamento com base na taxa de pacotes
		if packetRate > 20 {
			CONGESTION_THRESHOLD = 100
			CONGESTION_DELAY = 200 * time.Millisecond
			WINDOW_SIZE = 10

		} else if packetRate > 10 {
			CONGESTION_THRESHOLD = 75
			CONGESTION_DELAY = 400 * time.Millisecond
			WINDOW_SIZE = 15
		} else if packetRate > 5 {
			CONGESTION_THRESHOLD = 50
			CONGESTION_DELAY = 600 * time.Millisecond
			WINDOW_SIZE = 20
		} else {
			CONGESTION_THRESHOLD = 25
			CONGESTION_DELAY = 1000 * time.Millisecond
			WINDOW_SIZE = 50
		}
		
		color.Magenta("Adjusted CONGESTION_THRESHOLD to %d and CONGESTION_DELAY to %v based on packet rate: %.2f packets/sec\n", CONGESTION_THRESHOLD, CONGESTION_DELAY, packetRate)
		mutex.Unlock()
	}
}


