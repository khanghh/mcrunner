package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/khanghh/mcrunner/internal/mccmd"
	"github.com/khanghh/mcrunner/internal/sysmetrics"
	pb "github.com/khanghh/mcrunner/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MCRunnerService struct {
	pb.UnimplementedMCRunnerServer
	mcserver    *mccmd.MCServerCmd
	buffer      *ringBuffer
	consoleSubs map[grpc.ServerStreamingServer[pb.ConsoleMessage]]struct{}
	stateSubs   map[grpc.ServerStreamingServer[pb.ServerState]]struct{}
	done        chan struct{}
	mu          sync.Mutex
}

func (m *MCRunnerService) StartServer(ctx context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	if err := m.mcserver.Start(); err != nil {
		if errors.Is(err, mccmd.ErrAlreadyRunning) {
			return nil, status.Errorf(codes.Canceled, "Server is already running")
		}
		return nil, status.Errorf(codes.Internal, "Failed to start server: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (m *MCRunnerService) StopServer(ctx context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	if err := m.mcserver.Stop(); err != nil {
		if errors.Is(err, mccmd.ErrNotRunning) {
			return nil, status.Errorf(codes.Canceled, "Server is not running")
		}
		return nil, status.Errorf(codes.Internal, "Failed to stop server: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (m *MCRunnerService) KillServer(ctx context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	if err := m.mcserver.Kill(); err != nil {
		if errors.Is(err, mccmd.ErrNotRunning) {
			return nil, status.Errorf(codes.Canceled, "Server is not running")
		}
		return nil, status.Errorf(codes.Internal, "Failed to kill server: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (m *MCRunnerService) RestartServer(ctx context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	if err := m.mcserver.Stop(); err != nil {
		if errors.Is(err, mccmd.ErrNotRunning) {
			return nil, status.Errorf(codes.Canceled, "Server is not running")
		}
		return nil, status.Errorf(codes.Internal, "Failed to stop server: %v", err)
	}
	if err := m.mcserver.Start(); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to start server: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (m *MCRunnerService) SendCommand(ctx context.Context, cmdReq *pb.CommandRequest) (*emptypb.Empty, error) {
	if err := m.mcserver.SendCommand(cmdReq.Command); err != nil {
		if errors.Is(err, mccmd.ErrNotRunning) {
			return nil, status.Errorf(codes.Canceled, "Server is not running")
		}
		return nil, status.Errorf(codes.Internal, "Failed to send command: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (m *MCRunnerService) StreamConsole(stream grpc.BidiStreamingServer[pb.ConsoleMessage, pb.ConsoleMessage]) error {
	// register output subscriber
	m.mu.Lock()
	m.consoleSubs[stream] = struct{}{}
	m.mu.Unlock()

	// remove subscriber on exit
	defer func() {
		m.mu.Lock()
		delete(m.consoleSubs, stream)
		m.mu.Unlock()
	}()

	// read console input
	for {
		select {
		case <-m.done:
			return nil
		default:
		}

		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch payload := msg.Payload.(type) {
		case *pb.ConsoleMessage_PtyBuffer:
			if _, err := m.mcserver.Write(payload.PtyBuffer.Data); err != nil {
				errText := fmt.Sprintf("Failed to write console input: %v", err)
				fmt.Println(errText)
				if err := stream.Send(NewPtyErrorMessage(errText)); err != nil {
					return status.Errorf(codes.Unavailable, "Stream closed")
				}
			}
		case *pb.ConsoleMessage_PtyResize:
			rows, cols := int(payload.PtyResize.Rows), int(payload.PtyResize.Cols)
			if err := m.mcserver.ResizeWindow(rows, cols); err != nil {
				errText := fmt.Sprintf("Failed to resize PTY: %v", err)
				fmt.Println(errText)
				if err := stream.Send(NewPtyErrorMessage(errText)); err != nil {
					return status.Errorf(codes.Unavailable, "Stream closed")
				}
			}
		default:
			return status.Errorf(codes.InvalidArgument, "Unknown payload type")
		}
	}
}

func (m *MCRunnerService) StreamState(p0 *emptypb.Empty, stream grpc.ServerStreamingServer[pb.ServerState]) error {
	// add state subscriber
	m.mu.Lock()
	m.stateSubs[stream] = struct{}{}
	m.mu.Unlock()

	// remove subscriber on exit
	defer func() {
		m.mu.Lock()
		delete(m.stateSubs, stream)
		m.mu.Unlock()
	}()

	// Wait for client disconnect
	<-stream.Context().Done()
	return stream.Context().Err()
}

func (m *MCRunnerService) broadcastConsoleLoop() {
	broadcastCh := make(chan *pb.ConsoleMessage, 1)
	m.mcserver.OnStatusChanged(func(status mccmd.Status) {
		broadcastCh <- NewPtyStatusMessage(status)
	})
	go func() {
		stream := m.mcserver.OutputStream()
		buf := make([]byte, 4096)
		for {
			n, err := stream.Read(buf)
			if err != nil {
				return
			}

			data := make([]byte, n)
			copy(data, buf[:n])

			m.buffer.Write(buf[:n])

			select {
			case broadcastCh <- NewPtyBufferMessage(data):
			case <-m.done:
				return
			}
		}
	}()
	for {
		select {
		case msg := <-broadcastCh:
			m.mu.Lock()
			for stream := range m.consoleSubs {
				if err := stream.Send(msg); err != nil {
					log.Println("Failed to send PTY buffer message:", err)
				}
			}
			m.mu.Unlock()
		case <-m.done:
		}
	}
}

func (h *MCRunnerService) getServerState() *pb.ServerState {
	status := h.mcserver.GetStatus()
	var statusCode pb.Status
	switch status {
	case mccmd.StatusRunning:
		statusCode = pb.Status_STATUS_RUNNING
	case mccmd.StatusStopping:
		statusCode = pb.Status_STATUS_STOPPING
	case mccmd.StatusStopped:
		statusCode = pb.Status_STATUS_STOPPED
	default:
		statusCode = pb.Status_STATUS_UNKNOWN
	}
	serverState := &pb.ServerState{
		Status: statusCode,
	}

	usage, err := sysmetrics.GetResourceUsage()
	if err != nil {
		return serverState
	}
	serverState.MemoryUsage = usage.MemoryUsage
	serverState.MemoryLimit = usage.MemoryLimit
	serverState.CpuUsage = usage.CPUUsage
	serverState.CpuLimit = usage.CPULimit

	process := h.mcserver.GetProcess()
	if process == nil {
		return serverState
	}
	serverState.Pid = int32(process.Pid)

	if startTime := h.mcserver.GetStartTime(); startTime != nil {
		serverState.UptimeSec = uint64(time.Since(*startTime).Seconds())
	}
	return serverState
}

func (m *MCRunnerService) broadcastStateLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			state := m.getServerState()
			for stream := range m.stateSubs {
				if err := stream.Send(state); err != nil {
					log.Println("Failed to send server state message:", err)
				}
			}
			m.mu.Unlock()
		case <-m.done:
			return
		}
	}
}

func NewMCRunnerService(mcserver *mccmd.MCServerCmd) *MCRunnerService {
	svc := &MCRunnerService{
		mcserver:    mcserver,
		buffer:      newRingBuffer(1 << 20), // 1 MiB buffer
		consoleSubs: make(map[grpc.ServerStreamingServer[pb.ConsoleMessage]]struct{}),
		stateSubs:   make(map[grpc.ServerStreamingServer[pb.ServerState]]struct{}),
		done:        make(chan struct{}),
	}
	go svc.broadcastConsoleLoop()
	go svc.broadcastStateLoop()
	return svc
}
