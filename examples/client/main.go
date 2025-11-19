package main

import (
	"context"
	"fmt"

	"github.com/khanghh/mcrunner/pkg/api"
	pb "github.com/khanghh/mcrunner/pkg/proto"
)

func main() {
	url := "mcrunner://localhost:50051"
	client, err := api.NewMCRunnerGRPC(url)
	if err != nil {
		panic(err)
	}
	send := make(chan *pb.ConsoleMessage)
	receive := make(chan *pb.ConsoleMessage)
	go client.StreamConsole(context.Background(), send, receive)
	for msg := range receive {
		switch payload := msg.GetPayload().(type) {
		case *pb.ConsoleMessage_PtyBuffer:
			fmt.Print(string(payload.PtyBuffer.GetData()))
		case *pb.ConsoleMessage_PtyError:
		}
	}
}
