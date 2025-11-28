package params

import "time"

const (
	ServerBodyLimit              = 1 * 1024 * 1024 * 1024 // 1 MiB
	ServerIdleTimeout            = 30 * time.Second
	ServerReadTimeout            = 10 * time.Second
	ServerWriteTimeout           = 10 * time.Second
	ServerRunnerOutputBufferSize = 1024
	ServerRunnerErrorBufferSize  = 1024
	TTYBufferSize                = 4096
	WSClientQueueSize            = 128
)
