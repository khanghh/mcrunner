package websocket

import "errors"

var (
	ErrBadJSONMessage     = errors.New("bad JSON message")
	ErrUnknownMessageType = errors.New("unknown message type")
	ErrClientDisconnected = errors.New("client disconnected")
	ErrTopicNotFound      = errors.New("topic not found")
)
