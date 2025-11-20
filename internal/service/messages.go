package service

import (
	"github.com/khanghh/mcrunner/internal/mccmd"
	"github.com/khanghh/mcrunner/pkg/proto"
)

func NewPtyErrorMessage(message string) *proto.ConsoleMessage {
	return &proto.ConsoleMessage{
		Payload: &proto.ConsoleMessage_PtyError{
			PtyError: &proto.PtyError{
				Message: message,
			},
		},
	}
}

func NewPtyResizeMessage(rows, cols int) *proto.ConsoleMessage {
	return &proto.ConsoleMessage{
		Payload: &proto.ConsoleMessage_PtyResize{
			PtyResize: &proto.PtyResize{
				Rows: uint32(rows),
				Cols: uint32(cols),
			},
		},
	}
}

func NewPtyBufferMessage(output []byte) *proto.ConsoleMessage {
	return &proto.ConsoleMessage{
		Payload: &proto.ConsoleMessage_PtyBuffer{
			PtyBuffer: &proto.PtyBuffer{
				Data: output,
			},
		},
	}
}

func NewPtyStatusMessage(status mccmd.Status) *proto.ConsoleMessage {
	var pbStatus proto.Status
	switch status {
	case mccmd.StatusRunning:
		pbStatus = proto.Status_STATUS_RUNNING
	case mccmd.StatusStopping:
		pbStatus = proto.Status_STATUS_STOPPING
	case mccmd.StatusStopped:
		pbStatus = proto.Status_STATUS_STOPPED
	default:
		pbStatus = proto.Status_STATUS_UNKNOWN
	}
	return &proto.ConsoleMessage{
		Payload: &proto.ConsoleMessage_PtyStatus{
			PtyStatus: &proto.PtyStatus{
				Status: pbStatus,
			},
		},
	}
}

func NewServerStateMessage(state *proto.ServerState) *proto.ServerState {
	return &proto.ServerState{
		Status:      state.Status,
		Tps:         state.Tps,
		Pid:         state.Pid,
		MemoryUsage: state.MemoryUsage,
		MemoryLimit: state.MemoryLimit,
		CpuUsage:    state.CpuUsage,
		CpuLimit:    state.CpuLimit,
		UptimeSec:   state.UptimeSec,
	}
}
