package handlers

import (
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/pkg/gen"
)

func NewPTYOutputMessage(buf []byte) *gen.Message {
	return &gen.Message{
		Type: gen.MessageType_PTY_OUTPUT,
		Payload: &gen.Message_PtyBuffer{
			PtyBuffer: &gen.PtyBuffer{
				Data: buf,
			},
		},
	}
}

func NewServerStateMessage(status core.ServerStatus, pid int, tps float32, uptimeSec int64, serverUsage core.ServerUsage) *gen.Message {
	var statusCode gen.ServerStatus
	switch status {
	case core.StatusRunning:
		statusCode = gen.ServerStatus_RUNNING
	case core.StatusStopping:
		statusCode = gen.ServerStatus_STOPPING
	case core.StatusStopped:
		statusCode = gen.ServerStatus_STOPPED
	}

	return &gen.Message{
		Type: gen.MessageType_SERVER_STATE,
		Payload: &gen.Message_ServerState{
			ServerState: &gen.ServerState{
				Status:      statusCode,
				Pid:         int32(pid),
				Tps:         tps,
				UptimeSec:   uptimeSec,
				MemoryUsage: int64(serverUsage.MemoryUsage),
				MemoryLimit: int64(serverUsage.MemoryLimit),
				CpuUsage:    float64(serverUsage.CPUUsage),
				CpuLimit:    serverUsage.CPULimit,
			},
		},
	}
}
