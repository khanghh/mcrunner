package service

import (
	"context"

	"github.com/khanghh/mcrunner/internal/mccmd"
	pb "github.com/khanghh/mcrunner/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MCRunnerService struct {
	pb.UnimplementedMCRunnerServer
}

func (m *MCRunnerService) StartServer(p0 context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	panic("TODO: Implement")
}

func (m *MCRunnerService) StopServer(p0 context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	panic("TODO: Implement")
}

func (m *MCRunnerService) KillServer(p0 context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	panic("TODO: Implement")
}

func (m *MCRunnerService) RestartServer(p0 context.Context, p1 *emptypb.Empty) (*emptypb.Empty, error) {
	panic("TODO: Implement")
}

func (m *MCRunnerService) SendCommand(p0 context.Context, p1 *pb.CommandRequest) (*emptypb.Empty, error) {
	panic("TODO: Implement")
}

func (m *MCRunnerService) StreamConsole(p0 grpc.BidiStreamingServer[pb.ConsoleMessage, pb.ConsoleMessage]) error {
	panic("TODO: Implement")
}

func (m *MCRunnerService) StreamState(p0 *emptypb.Empty, p1 grpc.ServerStreamingServer[pb.ServerState]) error {
	panic("TODO: Implement")
}

func (m *MCRunnerService) mustEmbedUnimplementedMCRunnerServer() {
	panic("TODO: Implement")
}

func NewMCRunnerService(mcserver *mccmd.MCServerCmd) *MCRunnerService {
	return &MCRunnerService{}
}
