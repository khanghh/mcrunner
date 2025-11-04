package websocket

import (
	"context"
	"sync"

	"google.golang.org/protobuf/proto"
)

// HandlerFunc defines the function signature for a message handler.
type HandlerFunc func(ctx *Ctx, data []byte) error

type Ctx struct {
	context context.Context
	server  *Server
	client  *Client
	wg      sync.WaitGroup
}

func (c *Ctx) Context() context.Context {
	return c.context
}

func (c *Ctx) Client() *Client {
	return c.client
}

func (c *Ctx) SendError(msg string) error {
	c.server.sendError(c.client, msg)
	return nil
}

func (c *Ctx) SendMessage(msg *Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return c.client.Send(data)
}

func (c *Ctx) Disconnect(msg string) error {
	c.server.sendError(c.client, msg)
	c.client.Close()
	return nil
}

func (c *Ctx) Done() <-chan struct{} {
	return c.context.Done()
}

func (c *Ctx) Subscribe(topicName string) error {
	return c.server.clientSubscribe(c.client, topicName)
}

func (c *Ctx) Unsubscribe(topicName string) error {
	return c.server.clientUnsubscribe(c.client, topicName)
}

type Handler interface {
	Handle(ctx *Ctx, msg *Message) error
	OnConnect(ctx *Ctx)
}

type BaseHandler struct {
}

func (h *BaseHandler) Handle(ctx *Ctx, message *Message) error {
	return nil
}

func (h *BaseHandler) OnConnect(ctx *Ctx) error {
	return nil
}
