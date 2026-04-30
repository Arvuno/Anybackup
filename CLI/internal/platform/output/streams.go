package output

import (
	"io"
	"os"
)

type Streams struct {
	Stdout io.Writer
	Stderr io.Writer
}

func DefaultStreams() Streams {
	return Streams{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}
