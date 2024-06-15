package main

// import (
// 	"betterudp/client"
// 	"betterudp/server"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/go-playground/assert/v2"
// )

// const (
// 	TOTAL_MESSAGES = 10
// 	MAX_RETRIES    = 10
// 	TIMEOUT        = 1 * time.Second
// 	PORT           = ":1234"
// 	CONNECT        = "127.0.0.1:1234"
// )

// func TestSendMessages(t *testing.T) {
// 	// Start server and client
// 	go server.Server(PORT)
// 	go client.Client(CONNECT)

// 	expectedAck := ""

// 	// Wait for server and client to start (optional)
// 	time.Sleep(1 * time.Second)

// 	// Create UDP connection for sending messages
// 	s, err := net.ResolveUDPAddr("udp4", CONNECT)
// 	if err != nil {
// 		t.Fatalf("Failed to resolve UDP address: %v", err)
// 	}
// 	c, err := net.DialUDP("udp4", nil, s)
// 	if err != nil {
// 		t.Fatalf("Failed to dial UDP: %v", err)
// 	}
// 	defer c.Close()

// 	sentMessages := 0
// 	receivedAcks := 0

// 	for i := 0; i < TOTAL_MESSAGES; i++ {
// 		// Construct message
// 		message := "Message " + strconv.Itoa(i+1)
		
// 		retries := 0
// 		for {
// 			// Send message
// 			_, err = c.Write([]byte(message))
// 			if err != nil {
// 				t.Fatalf("Failed to send message: %v", err)
// 			}

// 			// Wait for ACK
// 			ackBuf := make([]byte, 1024)
// 			c.SetReadDeadline(time.Now().Add(TIMEOUT))
// 			n, _, err := c.ReadFromUDP(ackBuf)''
// 			if err != nil {
// 				fmt.Printf("No ACK received for message: %s, retrying...\n", message)
// 				retries++
// 				if retries >= MAX_RETRIES {
// 					fmt.Printf("Maximum retries exceeded for message: %s\n", message)
// 					break
// 				}
// 				continue
// 			}

// 			ack := string(ackBuf[:n])
// 			expectedAck = strconv.Itoa(i + 1)
// 			if ack == expectedAck {
// 				fmt.Printf("Received ACK: %s\n", ack)
// 				receivedAcks++
				
// 				break
// 			} else {
// 				fmt.Printf("Received incorrect ACK: %s, expecting: %s, retrying...\n", ack, expectedAck)
// 				retries++
// 				if retries >= MAX_RETRIES {
// 					fmt.Printf("Maximum retries exceeded for message: %s\n", message)
// 					break
// 				}
// 			}
// 		}

// 		sentMessages++
		
// 	}
// 	assert.Equal(t, sentMessages, TOTAL_MESSAGES)
// 	assert.Equal(t, strconv.Itoa(receivedAcks), expectedAck)
	
// }


// // func shuffleMessages(totalMessages int) []string {
// // 	messages := make([]string, totalMessages)
// // 	for i := 0; i < totalMessages; i++ {
// // 		messages[i] = "Message " + strconv.Itoa(i+1)
// // 	}

// // 	// Shuffle messages (Fisher-Yates shuffle algorithm)
// // 	rand.Seed(time.Now().UnixNano())
// // 	rand.Shuffle(len(messages), func(i, j int) {
// // 		messages[i], messages[j] = messages[j], messages[i]
// // 	})

// // 	return messages
// // }


// func TestSendMessages1000(t *testing.T) {
// 	// Start server and client
// 	go server.Server(PORT)
// 	go client.Client(CONNECT)

// 	expectedAck := ""

// 	// Wait for server and client to start (optional)
// 	time.Sleep(1 * time.Second)

// 	// Create UDP connection for sending messages
// 	s, err := net.ResolveUDPAddr("udp4", CONNECT)
// 	if err != nil {
// 		t.Fatalf("Failed to resolve UDP address: %v", err)
// 	}
// 	c, err := net.DialUDP("udp4", nil, s)
// 	if err != nil {
// 		t.Fatalf("Failed to dial UDP: %v", err)
// 	}
// 	defer c.Close()

// 	sentMessages := 0
// 	receivedAcks := 0

// 	startTime := time.Now()

// 	for i := 1000; i > 0; i-- {
// 		message := "Message " + strconv.Itoa(i)

// 		retries := 0
// 		for {
// 			// Send message
// 			sendStartTime := time.Now()
// 			_, err = c.Write([]byte(message))
// 			if err != nil {
// 				t.Fatalf("Failed to send message: %v", err)
// 			}
// 			sendDuration := time.Since(sendStartTime)
// 			log.Print(sendDuration) 
// 			// Wait for ACK
// 			ackBuf := make([]byte, 1024)
// 			c.SetReadDeadline(time.Now().Add(TIMEOUT))
// 			n, _, err := c.ReadFromUDP(ackBuf)
// 			if err != nil {
// 				fmt.Printf("No ACK received for message: %s, retrying...\n", message)
// 				retries++
// 				if retries >= MAX_RETRIES {
// 					fmt.Printf("Maximum retries exceeded for message: %s\n", message)
// 					break
// 				}
// 				continue
// 			}

// 			ack := string(ackBuf[:n])
// 			expectedAck = strconv.Itoa(i)
// 			if ack == expectedAck {
// 				fmt.Printf("Received ACK: %s\n", ack)
// 				receivedAcks++
// 				break
// 			} else {
// 				fmt.Printf("Received incorrect ACK: %s, expecting: %s, retrying...\n", ack, expectedAck)
// 				retries++
// 				if retries >= MAX_RETRIES {
// 					fmt.Printf("Maximum retries exceeded for message: %s\n", message)
// 					break
// 				}
// 			}
// 		}
		
// 		sentMessages++
// 	}

// 	// Calcular o tempo total decorrido
// 	totalDuration := time.Since(startTime)
	
// 	// Logs para tempo total e médio por mensagem
// 	fmt.Printf("Total elapsed time: %s\n", totalDuration)
// 	fmt.Printf("Average time per message: %s\n", totalDuration/time.Duration(1000))

// 	// Assertivas para verificar o número de mensagens enviadas e os ACKs recebidos
// 	assert.Equal(t, sentMessages, 1000)
// 	assert.Equal(t, strconv.Itoa(receivedAcks), expectedAck)
// }
