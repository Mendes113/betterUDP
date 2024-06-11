package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "time"
)

const (
    WINDOW_SIZE = 50
    TOTAL_MESSAGES = 1000
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

    fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
    defer c.Close()

    for i := 0; i < TOTAL_MESSAGES; i++ {
        // Construct message
        message := "Message " + strconv.Itoa(i+1)

        // Send message
        _, err = c.Write([]byte(message))
        if err != nil {
            fmt.Println(err)
            return
        }

        // Wait for a short duration to simulate real-time communication
        time.Sleep(time.Millisecond * 100)
    }

    fmt.Printf("Sent %d messages\n", TOTAL_MESSAGES)
}
