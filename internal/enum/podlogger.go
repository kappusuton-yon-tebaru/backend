package enum

type PodLoggerMode string

const (
	PodLoggerModeStdout  = "log"
	PodLoggerModeMongoDb = "mongodb"
)

