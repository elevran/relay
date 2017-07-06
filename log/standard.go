package log

import (
	"log"
	"os"
)

// NewStdlibLogger returns a logger implementation backed by the standard library log package
func NewStdlibLogger(prefix string) Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lshortfile)
}
