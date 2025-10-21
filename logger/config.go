package logger

type Config struct {
	Level           string // debug, info, warn, error
	Encoding        string // json or console
	Service         string
	Environment     string
	AsyncBufferSize int // number of log entries in buffer
	BatchSize       int // number of logs per batch flush
}
