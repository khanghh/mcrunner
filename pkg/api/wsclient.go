package api

import (
	"io"
	"log"

	fws "github.com/fasthttp/websocket"
)

// WebsocketClient provides a reader interface for WebSocket streams
type WebsocketClient struct {
	conn   *fws.Conn
	ch     chan []byte
	stopCh chan struct{}
}

// NewWebsocketClient creates a new WebSocket reader for the given URL
func DialWebsocket(wsURL string) (*WebsocketClient, error) {
	conn, _, err := fws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	wr := &WebsocketClient{
		conn:   conn,
		ch:     make(chan []byte, 128), // Buffer size matching params.WSClientQueueSize
		stopCh: make(chan struct{}),
	}

	go wr.readLoop()

	return wr, nil
}

// readLoop continuously reads from the WebSocket connection
func (wc *WebsocketClient) readLoop() {
	defer wc.conn.Close()

	for {
		select {
		case <-wc.stopCh:
			return
		default:
			_, data, err := wc.conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}
			select {
			case wc.ch <- data:
			case <-wc.stopCh:
				return
			default:
				// Drop message if channel is full (non-blocking)
			}
		}
	}
}

// Read implements io.Reader
func (wc *WebsocketClient) Read(p []byte) (int, error) {
	select {
	case data := <-wc.ch:
		n := copy(p, data)
		return n, nil
	case <-wc.stopCh:
		return 0, io.EOF
	}
}

func (wc *WebsocketClient) Write(buf []byte) (int, error) {
	err := wc.conn.WriteMessage(fws.BinaryMessage, buf)
	if err != nil {
		return 0, err
	}
	return len(buf), nil
}

// Close closes the WebSocket connection
func (wc *WebsocketClient) Close() error {
	select {
	case <-wc.stopCh:
		return nil
	default:
		close(wc.stopCh)
		return wc.conn.Close()
	}
}
