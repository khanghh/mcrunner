package handlers

import (
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/pkg/gen"
)

func NewServerStatusMessage(state core.ServerState) *gen.Message {
	var stateCode gen.ServerState
	switch state {
	case core.StateRunning:
		stateCode = gen.ServerState_RUNNING
	case core.StateStopping:
		stateCode = gen.ServerState_STOPPING
	case core.StateStopped:
		stateCode = gen.ServerState_STOPPED
	}

	return &gen.Message{
		Type: gen.MessageType_SERVER_STATUS,
		Payload: &gen.Message_ServerStatus{
			ServerStatus: &gen.ServerStatus{
				State: stateCode,
			},
		},
	}
}
