package main

import (
	"net/http"
	"os"
	"time"

	"github.com/anybackup-ai/Anybackup/CLI/cmd"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

var cliVersion = "dev"

func main() {
	streams := output.DefaultStreams()
	root := cmd.NewRootCommand(cmd.Dependencies{
		Streams: streams,
		Version: cmd.BuildVersionInfo(cliVersion),
		Console: cmd.NewDefaultExecutor(&http.Client{Timeout: 30 * time.Second}),
	})

	err := root.Execute()
	if err != nil {
		_ = output.WriteError(streams, err)
	}
	os.Exit(clierrors.ExitCodeOf(err))
}
