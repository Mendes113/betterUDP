package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	// "time"

	// "betterudp/client"
	// "betterudp/server"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const CSV_FILE_PATH = "/home/mendes/Documents/Github/betterUDP/execution_times.csv"

var mutex sync.Mutex

func main() {
	// // Iniciar temporizador para medir o tempo do servidor
	// serverStartTime := time.Now()

	// // Iniciar servidor em uma goroutine
	// go server.Server(":1234")

	// // Aguardar um breve momento para garantir que o servidor esteja pronto
	// time.Sleep(100 * time.Millisecond)

	// // Medir tempo de execução do cliente
	// clientStartTime := time.Now()

	// // Executar o cliente
	// client.Client("127.0.0.1:1234")

	// // Registrar tempo total de execução do cliente
	// clientEndTime := time.Now()
	// clientTotalTime := clientEndTime.Sub(clientStartTime)

	// // Encerrar o servidor
	// serverEndTime := time.Now()
	// serverTotalTime := serverEndTime.Sub(serverStartTime)

	// fmt.Printf("Tempo total de execução do servidor: %v\n", serverTotalTime)
	// fmt.Printf("Tempo total de execução do cliente: %v\n", clientTotalTime)

	// // Salvar tempos de execução no CSV
	// err := LogElapsedTime(serverTotalTime.Milliseconds(), clientTotalTime.Milliseconds())
	// if err != nil {
	// 	fmt.Println("Erro ao salvar tempos de execução no CSV:", err)
	// }

	err := GenerateExecutionTimeChart("/home/mendes/Documents/Github/betterUDP/execution_times.csv", "output.png")
	if err != nil {
		fmt.Println("Erro:", err)
	}
}

func LogElapsedTime(serverTimeMs, clientTimeMs int64) error {
	// Abrir o arquivo CSV no modo de apêndice
	file, err := os.OpenFile(CSV_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo CSV: %w", err)
	}
	defer file.Close()

	// Criar um escritor CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escrever os tempos de execução no CSV
	mutex.Lock()
	defer mutex.Unlock()

	

	if err := writer.Write([]string{strconv.FormatInt(serverTimeMs, 10), strconv.FormatInt(clientTimeMs, 10)}); err != nil {
		return fmt.Errorf("erro ao escrever tempos de execução no CSV: %w", err)
	}

	return nil
}

func GenerateExecutionTimeChart(csvFilePath, outputImagePath string) error {
	// Abrir o arquivo CSV
	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo CSV: %w", err)
	}
	defer file.Close()

	// Ler os dados do CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo CSV: %w", err)
	}

	// Listas para armazenar os tempos de execução
	var serverTimes []float64
	// var clientTimes []float64
	var serverTCPTimes []float64
	// var clientTCPTimes []float64
	var serverUDPTimes []float64
	// var clientUDPTimes []float64

	// Iterar sobre os registros do CSV
	for _, record := range records {
		// Verificar se a linha possui exatamente 4 valores
		if len(record) != 6 {
			continue // Ignorar esta linha se não tiver exatamente 4 valores
		}

		// Converter e armazenar os tempos de execução
		serverTime, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return fmt.Errorf("erro ao converter tempo de execução do servidor: %w", err)
		}
		// clientTime, err := strconv.ParseFloat(record[1], 64)
		// if err != nil {
		// 	return fmt.Errorf("erro ao converter tempo de execução do cliente: %w", err)
		// }
		serverTCPTime, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return fmt.Errorf("erro ao converter tempo de execução do servidor TCP: %w", err)
		}
		// clientTCPTime, err := strconv.ParseFloat(record[3], 64)
		// if err != nil {
		// 	return fmt.Errorf("erro ao converter tempo de execução do cliente TCP: %w", err)
		// }
		// clientUDPTime, err := strconv.ParseFloat(record[4], 64)
		// if err != nil {
		// 	return fmt.Errorf("erro ao converter tempo de execução do cliente UDP: %w", err)
		// }
		serverUDPTime, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return fmt.Errorf("erro ao converter tempo de execução do servidor UDP: %w", err)
		}


		// Adicionar os tempos à lista correspondente
		serverTimes = append(serverTimes, serverTime)
		// clientTimes = append(clientTimes, clientTime)
		serverTCPTimes = append(serverTCPTimes, serverTCPTime)
		// clientTCPTimes = append(clientTCPTimes, clientTCPTime)
		// clientUDPTimes = append(clientUDPTimes, clientUDPTime)
		serverUDPTimes = append(serverUDPTimes, serverUDPTime)
	}


	// Calcular média dos tempos de execução
	avgServerTime := calculateAverage(serverTimes)
	// avgClientTime := calculateAverage(clientTimes)
	avgServerTCPTime := calculateAverage(serverTCPTimes)
	// avgClientTCPTime := calculateAverage(clientTCPTimes)
	// avgClientUDPTime := calculateAverage(clientUDPTimes)
	avgServerUDPTime := calculateAverage(serverUDPTimes)

	fmt.Printf("Média de tempo de execução do servidor: %.2f ms\n", avgServerTime)
	// fmt.Printf("Média de tempo de execução do cliente: %.2f ms\n", avgClientTime)
	fmt.Printf("Média de tempo de execução do servidor TCP: %.2f ms\n", avgServerTCPTime)
	// fmt.Printf("Média de tempo de execução do cliente TCP: %.2f ms\n", avgClientTCPTime)
	fmt.Printf("Média de tempo de execução do servidor UDP: %.2f ms\n", avgServerUDPTime)
	// fmt.Printf("Média de tempo de execução do cliente UDP: %.2f ms\n", avgClientUDPTime)

	// Criar plot
	p := plot.New()

	// Definir título e rótulos dos eixos
	p.Title.Text = "Tempo Médio de Execução"
	p.X.Label.Text = "Execuções"
	p.Y.Label.Text = "Tempo (ms)"

	// Criar pontos para plotar
	pointsServer := make(plotter.XYs, len(serverTimes))
	// pointsClient := make(plotter.XYs, len(clientTimes))
	pointsServerTCP := make(plotter.XYs, len(serverTCPTimes))
	// pointsClientTCP := make(plotter.XYs, len(clientTCPTimes))
	// pointsClientUDP := make(plotter.XYs, len(clientUDPTimes))
	pointsServerUDP := make(plotter.XYs, len(serverUDPTimes))

	for i := range pointsServer {
		pointsServer[i].X = float64(i + 1)
		pointsServer[i].Y = serverTimes[i]
		// pointsClient[i].X = float64(i + 1)
		// pointsClient[i].Y = clientTimes[i]
		pointsServerTCP[i].X = float64(i + 1)
		pointsServerTCP[i].Y = serverTCPTimes[i]
		// pointsClientTCP[i].X = float64(i + 1)
		// pointsClientTCP[i].Y = clientTCPTimes[i]
		// pointsClientUDP[i].X = float64(i + 1)
		pointsServerUDP[i].Y = serverUDPTimes[i]
	}

	// Adicionar pontos ao plot
	err = plotutil.AddLinePoints(p,
		"Betterudp Server", pointsServer,
		// "Betterudp Client", pointsClient,
		"Servidor TCP", pointsServerTCP,
		// "Cliente TCP", pointsClientTCP,
		// "Cliente UDP", pointsClientUDP,
		"Servidor UDP", pointsServerUDP,
	
	)
	if err != nil {
		return fmt.Errorf("erro ao adicionar pontos ao plot: %w", err)
	}

	// Salvar plot como imagem PNG
	if err := p.Save(8*vg.Inch, 4*vg.Inch, outputImagePath); err != nil {
		return fmt.Errorf("erro ao salvar plot como imagem PNG: %w", err)
	}

	return nil
}

	
	// Função para calcular a média de um slice de float64
	func calculateAverage(times []float64) float64 {
		sum := 0.0
		for _, t := range times {
			sum += t
		}
		return sum / float64(len(times))
	}