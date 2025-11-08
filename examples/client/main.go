package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/fasthttp/websocket"
)

func main() {
	var addr = flag.String("addr", "localhost:3000", "http service address")
	flag.Parse()

	// Parse the URL
	u, err := url.Parse(fmt.Sprintf("ws://%s/ws", *addr))
	if err != nil {
		log.Fatal("Failed to parse URL:", err)
	}

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket:", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", u.String())
	fmt.Println("Type commands to send to the server. Ctrl+C to exit.")

	// Handle interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel to signal when done
	done := make(chan struct{})

	// Goroutine to read from WebSocket and print to stdout
	go func() {
		defer close(done)
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					fmt.Println("\nServer disconnected")
					return
				}
				log.Println("Read error:", err)
				return
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				fmt.Print(string(p))
			}
		}
	}()

	// Read from stdin and send to WebSocket
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("\nExiting...")
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			return
		default:
			if scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				err := conn.WriteMessage(websocket.TextMessage, []byte(line+"\n"))
				if err != nil {
					log.Println("Send error:", err)
					return
				}
			} else {
				return
			}
		}
	}
}
