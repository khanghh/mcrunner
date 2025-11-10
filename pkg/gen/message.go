package gen

func NewPTYBufferMessage(buf []byte) *Message {
	return &Message{
		Type: MessageType_PTY_BUFFER,
		Payload: &Message_PtyBuffer{
			PtyBuffer: &PtyBuffer{
				Data: buf,
			},
		},
	}
}
