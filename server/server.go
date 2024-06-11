package main

import (
    "fmt"
    "math/rand"
    "net"
    "os"
    "strconv"
    "strings"
    "time"
)

const (
    WINDOW_SIZE = 50
)

var (
    lastAckedSeqNum int
    nextSeqNum      int
    congestion      bool
)

func random(min, max int) int {
    return rand.Intn(max-min) + min
}

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
    rand.Seed(time.Now().Unix())

    for {
        n, addr, err := connection.ReadFromUDP(buffer)
        fmt.Print("-> ", string(buffer[0:n-1]))

        if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
            fmt.Println("Exiting UDP server!")
            return
        }

        data := []byte(strconv.Itoa(random(1, 1001)))
        fmt.Printf("data: %s\n", string(data))

        // Simulating congestion
        if congestion {
            time.Sleep(time.Millisecond * time.Duration(random(50, 200)))
            congestion = false
        }

        // Send data if within window size
        if nextSeqNum-lastAckedSeqNum < WINDOW_SIZE {
            seqData := append([]byte(strconv.Itoa(nextSeqNum)), data...)
            _, err = connection.WriteToUDP(seqData, addr)
            if err != nil {
                fmt.Println(err)
                return
            }
            nextSeqNum++
        } else {
            congestion = true
            fmt.Println("Congestion detected! Waiting...")
        }
    }
}
