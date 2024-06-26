package client
import (
	"fmt"
	"net"
	"strconv"
	"time"
    "github.com/fatih/color"
)
const (
	TOTAL_MESSAGES = 100000
	MAX_RETRIES    = 100000
	TIMEOUT        = 1 * time.Second
)
func sendMessage(c *net.UDPConn, message string, seqNum int) error {
    retries := 0
    for {
        now := time.Now()
        // Send message
        _, err := c.Write([]byte(message))
        if err != nil {
            return err
        }
        // Calcula tempo decorrido para enviar a mensagem
        sendElapsed := time.Since(now)
		fmt.Printf("Sent message '%s', took %v\n", message, sendElapsed)
        // cor da mensagem
        color.Cyan("Sent message '%s', took %v\n", message, sendElapsed)
        // ack buffer
        ackBuf := make([]byte, 1024)
        // Set read deadline
        c.SetReadDeadline(time.Now().Add(TIMEOUT))
        // le o ack
        n, _, err := c.ReadFromUDP(ackBuf)
        if err != nil {
            fmt.Printf("No ACK received for message: %s, retrying...\n", message)
            retries++
            // Se o número de tentativas exceder o limite, retorna um erro
            if retries >= MAX_RETRIES {
                fmt.Printf("Maximum retries exceeded for message: %s, unable to send further messages.\n", message)
                return fmt.Errorf("maximum retries exceeded")
            }
            continue
        }
        // converte o ack para string
        ack := string(ackBuf[:n])
        // verifica se o ack é o esperado
        expectedAck := strconv.Itoa(seqNum)
        if ack == expectedAck {
            fmt.Printf("Received ACK: %s\n", ack)
            break
        } else {
            fmt.Printf("Received incorrect ACK: %s, expecting: %s, retrying...\n", ack, expectedAck)
            retries++
            if retries >= MAX_RETRIES {
                fmt.Printf("Maximum retries exceeded for message: %s, unable to send further messages.\n", message)
                return fmt.Errorf("maximum retries exceeded")
            }
        }
    }
    return nil
}
func Client(CONNECT string) {
    // Resolve o endereço do servidor
	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
    // Abre a conexão UDP
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	start := time.Now() // Marca o tempo de início do envio das mensagens
    // Envia as mensagens
	for i := 0; i < TOTAL_MESSAGES; i++ {
		message := "Message " + strconv.Itoa(i+1)
        // Envia a mensagem
		err := sendMessage(c, message, i+1)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			return
		}
	}
	elapsed := time.Since(start) // Calcula o tempo total decorrido
	fmt.Printf("Sent %d messages in %v\n", TOTAL_MESSAGES, elapsed)
}
