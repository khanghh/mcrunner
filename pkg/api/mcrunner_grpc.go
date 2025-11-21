package api

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/khanghh/mcrunner/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	resolver.Register(&mcrunnerBuilder{})
}

type ConsoleMessageHandler func(msg *pb.ConsoleMessage)

type MCRunnerGRPC struct {
	conn *grpc.ClientConn
	cl   pb.MCRunnerClient
}

func (c *MCRunnerGRPC) StartServer(ctx context.Context) error {
	_, err := c.cl.StartServer(ctx, &emptypb.Empty{})
	return err
}

func (c *MCRunnerGRPC) StopServer(ctx context.Context) error {
	_, err := c.cl.StopServer(ctx, &emptypb.Empty{})
	return err
}
func (c *MCRunnerGRPC) KillServer(ctx context.Context) error {
	_, err := c.cl.KillServer(ctx, &emptypb.Empty{})
	return err
}

func (c *MCRunnerGRPC) RestartServer(ctx context.Context) error {
	_, err := c.cl.RestartServer(ctx, &emptypb.Empty{})
	return err
}

func (c *MCRunnerGRPC) SendCommand(ctx context.Context, cmd string) error {
	_, err := c.cl.SendCommand(ctx, &pb.CommandRequest{
		Command: cmd,
	})
	return err
}

func (c *MCRunnerGRPC) GetState(ctx context.Context) (*pb.ServerState, error) {
	return c.cl.GetState(ctx, &emptypb.Empty{})
}

func (c *MCRunnerGRPC) handleStreamConsole(ctx context.Context, stream pb.MCRunner_StreamConsoleClient, send <-chan *pb.ConsoleMessage, receive chan<- *pb.ConsoleMessage) error {
	errChan := make(chan error, 2)

	// Send goroutine
	go func() {
		for {
			select {
			case msg, ok := <-send:
				if !ok {
					errChan <- stream.CloseSend()
					return
				}
				if err := stream.Send(msg); err != nil {
					errChan <- err
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Receive loop
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					errChan <- nil
					return
				}
				errChan <- err
			}

			select {
			case receive <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (c *MCRunnerGRPC) StreamConsole(ctx context.Context, send <-chan *pb.ConsoleMessage, receive chan<- *pb.ConsoleMessage) error {
	defer close(receive)
	for {
		// open stream
		stream, err := c.cl.StreamConsole(ctx)
		if err != nil {
			log.Println("Failed to open console stream:", err)
		} else {
			// handle bidirectional stream
			streamCtx, cancel := context.WithCancel(ctx)
			if err := c.handleStreamConsole(streamCtx, stream, send, receive); err != nil {
				log.Println("Console stream closed:", err)
			}
			cancel()
		}

		// Wait 1 second before reconnecting
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *MCRunnerGRPC) StreamState(ctx context.Context, receive chan<- *pb.ServerState) error {
	defer close(receive)
	for {
		stream, err := c.cl.StreamState(ctx, &emptypb.Empty{})
		if err != nil {
			log.Println("Failed to open state stream:", err)
		} else {
			// Receive loop
			for {
				state, err := stream.Recv()
				if err != nil {
					log.Println("State stream closed:", err)
					break
				}

				select {
				case receive <- state:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		// Wait 1 second before reconnecting
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *MCRunnerGRPC) Close() error {
	return c.conn.Close()
}

func NewMCRunnerGRPC(addr string) (*MCRunnerGRPC, error) {
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
	return &MCRunnerGRPC{
		conn: conn,
		cl:   pb.NewMCRunnerClient(conn),
	}, nil
}
