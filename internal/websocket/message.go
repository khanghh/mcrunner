package websocket

func NewPTYBufferMessage(sid string, data []byte) *Message {
	return &Message{
		Type: MessageType_PTY_BUFFER,
		Payload: &Message_PtyBuffer{
			PtyBuffer: &PtyBuffer{
				SessionId: sid,
				Data:      data,
			},
		},
	}
}
