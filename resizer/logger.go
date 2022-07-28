package resizer

import "log"

func newLogger() *log.Logger {
	logger := log.Default()
	logger.SetPrefix("[image-resizer] ")
	return logger
}
