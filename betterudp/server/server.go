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
	packetCount           int
	lastPacketTime        time.Time
	CONGESTION_THRESHOLD  = 100                        // Congestion threshold initial
	CONGESTION_DELAY      = 100 * time.Millisecond     // Congestion delay initial
	CSV_FILE_PATH         = "./times.csv"
	startTime, endTime    time.Time
	WINDOW_SIZE           = 2 // Initial congestion window size
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
	// go routine é similar a uma thread
	// entretanto em go uma go routine é mais leve que uma thread
	// diversas go routines podem ser executadas em uma única thread
	go AdjustCongestionParameters()
	// Inicia a goroutine para calcular a taxa de pacotes
	go CalculatePacketRate()
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
	// mutex serve para garantir que apenas uma go routine execute o código por vez
	mutex.Lock()
	// defer garante que a função seja chamada no final da execução da função atual
	defer mutex.Unlock()
	
	// Incrementa o contador de pacotes
	packetCount++

	// Extrai e converte o número de sequência da mensagem recebida
    // A mensagem é esperada ter o formato "algum_texto <numero_de_sequencia>"
    // Divide a mensagem pelo espaço e tenta converter o segundo elemento para inteiro
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
		// Converte o próximo número de sequência esperado para string
		expectedSeqNum := strconv.Itoa(baseSeqNum)
		// Converte a string do número de sequência esperado para um slice de bytes
		ack := []byte(expectedSeqNum)
		// Envia o ACK para o endereço de origem
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

func CalculatePacketRate() {
	for {
		time.Sleep(1 * time.Second)
		mutex.Lock()
		packetRate = float64(packetCount)
		packetCount = 0 // Reseta o contador após calcular a taxa
		mutex.Unlock()
	}
}

func AdjustCongestionParameters() {
	// Ajusta os parâmetros de congestionamento
	for {
		// ajusta a cada segundo
		time.Sleep(1 * time.Second) 
		// mutex serve para garantir que apenas uma go routine execute o código por vez
		mutex.Lock()
		// ajusta a janela de congestionamento com base na taxa de pacotes
		if packetRate > 25 {
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
		color.Magenta("Adjusted CONGESTION_THRESHOLD to %d and CONGESTION_DELAY %v and Window Size to %d based on packet rate: %.2f packets/sec\n", CONGESTION_THRESHOLD, CONGESTION_DELAY, WINDOW_SIZE, packetRate)
		mutex.Unlock()
	}
}
