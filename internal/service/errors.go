package service

import (
	"errors"

	"github.com/khanghh/mcrunner/internal/mccmd"
	pb "github.com/khanghh/mcrunner/pkg/proto"
)

func mapMCCmdError(err error) *pb.ConsoleMessage {
	if errors.Is(err, mccmd.ErrNotRunning) {
		return newPtyErrorMessage("NOT_RUNNING", err.Error())
	}
	if errors.Is(err, mccmd.ErrAlreadyRunning) {
		return newPtyErrorMessage("ALREADY_RUNNING", err.Error())
	}
	return newPtyErrorMessage("UNKNOWN", err.Error())
}

func newPtyErrorMessage(code string, message string) *pb.ConsoleMessage {
	return &pb.ConsoleMessage{
		Payload: &pb.ConsoleMessage_PtyError{
			PtyError: &pb.PtyError{
				Code:    code,
				Message: message,
			},
		},
	}
}
