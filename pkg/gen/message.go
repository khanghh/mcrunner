package gen

func NewPTYOutputMessage(buf []byte) *Message {
	return &Message{
		Type: MessageType_PTY_OUTPUT,
		Payload: &Message_PtyBuffer{
			PtyBuffer: &PtyBuffer{
				Data: buf,
			},
		},
	}
}
