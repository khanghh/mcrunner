package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	gws "github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	wsproto "github.com/khanghh/mcrunner/internal/websocket"
)

func main() {
	url := flag.String("url", "ws://127.0.0.1:3000/ws", "WebSocket URL, e.g. ws://127.0.0.1:3000/ws")
	flag.Parse()

	log.Printf("Connecting to %s", *url)
	conn, _, err := gws.DefaultDialer.Dial(*url, nil)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	// Reader goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			mt, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				return
			}
			if mt != gws.BinaryMessage && mt != gws.TextMessage {
				continue
			}

			var msg wsproto.Message
			if err := proto.Unmarshal(data, &msg); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}
			printMessage(&msg)
		}
	}()

	// Graceful shutdown on Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-done:
		log.Println("server closed connection")
	case <-interrupt:
		log.Println("interrupt: closing")
		_ = conn.WriteMessage(gws.CloseMessage, gws.FormatCloseMessage(gws.CloseNormalClosure, "bye"))
		select {
		case <-done:
		case <-time.After(2 * time.Second):
		}
	}
}

func printMessage(m *wsproto.Message) {
	switch m.GetType() {
	case wsproto.MessageType_ERROR:
		log.Printf("[ERROR] %s", m.GetError())

	case wsproto.MessageType_PTY_BUFFER:
		if pb := m.GetPtyBuffer(); pb != nil {
			fmt.Print(string(pb.GetData()))
			// log.Printf("[PTY_BUFFER] session=%s bytes=%d data=%q",
			// 	pb.GetSessionId(), len(pb.GetData()), string(pb.GetData()))
		}

	case wsproto.MessageType_PTY_INPUT:
		if pi := m.GetPtyInput(); pi != nil {
			log.Printf("[PTY_INPUT] session=%s bytes=%d data=%q",
				pi.GetSessionId(), len(pi.GetData()), string(pi.GetData()))
		}

	case wsproto.MessageType_PTY_RESIZE:
		if pr := m.GetPtyResize(); pr != nil {
			log.Printf("[PTY_RESIZE] session=%s cols=%d rows=%d",
				pr.GetSessionId(), pr.GetCols(), pr.GetRows())
		}

	default:
		log.Printf("[UNKNOWN] type=%v", m.GetType())
	}
}
