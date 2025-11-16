package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/khanghh/mcrunner/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MCRunnerClient struct {
	conn  *grpc.ClientConn
	cl    pb.MCRunnerClient
	alive bool
}

func (c *MCRunnerClient) IsAlive() bool {
	return c.alive
}

func (c *MCRunnerClient) StartServer(ctx context.Context) error {
	_, err := c.cl.StartServer(ctx, &emptypb.Empty{})
	return err
}

func (c *MCRunnerClient) StopServer(ctx context.Context) error {
	_, err := c.cl.StopServer(ctx, &emptypb.Empty{})
	return err
}
func (c *MCRunnerClient) KillServer(ctx context.Context) error {
	_, err := c.cl.KillServer(ctx, &emptypb.Empty{})
	return err
}

func (c *MCRunnerClient) RestartServer(ctx context.Context) error {
	_, err := c.cl.RestartServer(ctx, &emptypb.Empty{})
	return err
}

func (c *MCRunnerClient) SendCommand(ctx context.Context, cmd string) error {
	_, err := c.cl.SendCommand(ctx, &pb.CommandRequest{
		Command: cmd,
	})
	return err
}

func (c *MCRunnerClient) handleStreamConsole(ctx context.Context, stream pb.MCRunner_StreamConsoleClient, send <-chan *pb.ConsoleMessage, receive chan<- *pb.ConsoleMessage) error {
	errChan := make(chan error, 1)

	// Send goroutine
	go func() {
		defer close(errChan)
		fmt.Println("send loop started")
		for {
			select {
			case msg, ok := <-send:
				if !ok {
					errChan <- stream.CloseSend()
					fmt.Println("send closed 1")
					return
				}
				if err := stream.Send(msg); err != nil {
					errChan <- err
					fmt.Println("send closed 2")
					return
				}
			case <-ctx.Done():
				fmt.Println("send closed 3")
				return
			}
		}
	}()

	// Receive loop
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		select {
		case receive <- msg:
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *MCRunnerClient) StreamConsole(ctx context.Context, send <-chan *pb.ConsoleMessage, receive chan<- *pb.ConsoleMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		stream, err := c.cl.StreamConsole(ctx)
		if err != nil {
			log.Println("Failed to create stream, retrying:", err)
			time.Sleep(1 * time.Second) // backoff here if you want
			continue
		}

		err = c.handleStreamConsole(ctx, stream, send, receive)
		if err != nil {
			log.Println("Stream error, reconnecting:", err)
			time.Sleep(1 * time.Second) // optional backoff
			continue
		}

		return nil
	}
}

func (c *MCRunnerClient) Close() error {
	return c.conn.Close()
}

func NewMCRunnerClient(addr string) (*MCRunnerClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Second,
				Multiplier: 1,
				MaxDelay:   time.Second,
			},
			MinConnectTimeout: time.Second,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second, // ping interval
			Timeout: 10 * time.Second, // ping ack timeout
		}),
	)
	if err != nil {
		return nil, err
	}
	return &MCRunnerClient{
		conn: conn,
		cl:   pb.NewMCRunnerClient(conn),
	}, nil
}
