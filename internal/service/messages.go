package service

import "github.com/khanghh/mcrunner/pkg/proto"

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
