package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	TOTAL_MESSAGES = 100
	MAX_RETRIES    = 100    
	TIMEOUT        = 1 * time.Second
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer c.Close()

	for i := 0; i < TOTAL_MESSAGES; i++ {
		// Construct message
		message := "Message " + strconv.Itoa(i+1)

		retries := 0
		for {
			// Send message
			_, err = c.Write([]byte(message))
			if err != nil {
				fmt.Println(err)
				return
			}

			// Wait for ACK
			ackBuf := make([]byte, 1024)
			c.SetReadDeadline(time.Now().Add(TIMEOUT))
			n, _, err := c.ReadFromUDP(ackBuf)
			if err != nil {
				fmt.Printf("No ACK received for message: %s, retrying...\n", message)
				retries++
				if retries >= MAX_RETRIES {
					fmt.Printf("Maximum retries exceeded for message: %s, unable to send further messages.\n", message)
					return
				}
				continue
			}

			ack := string(ackBuf[:n])
			expectedAck := strconv.Itoa(i + 1)
			if ack == expectedAck {
				fmt.Printf("Received ACK: %s\n", ack)
				break
			} else {
				fmt.Printf("Received incorrect ACK: %s, expecting: %s, retrying...\n", ack, expectedAck)
				retries++
				if retries >= MAX_RETRIES {
					fmt.Printf("Maximum retries exceeded for message: %s, unable to send further messages.\n", message)
					return
				}
			}
		}

		// Wait for a short duration to simulate real-time communication
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("Sent %d messages\n", TOTAL_MESSAGES)
}
