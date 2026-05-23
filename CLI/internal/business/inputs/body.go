package inputs

import (
	"encoding/json"
	"os"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

func LoadBody(inline string, bodyFile string) ([]byte, error) {
	switch {
	case inline != "" && bodyFile != "":
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "-d and --body-file are mutually exclusive")
	case inline == "" && bodyFile == "":
		return nil, clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
	case inline != "":
		if !json.Valid([]byte(inline)) {
			return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "inline body is not valid json")
		}
		return []byte(inline), nil
	default:
		body, err := os.ReadFile(bodyFile)
		if err != nil {
			return nil, clierrors.Wrap(clierrors.CodeBodyFileReadFailed, clierrors.ExitInvalidArgs, "failed to read --body-file", err)
		}
		if !json.Valid(body) {
			return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "body file is not valid json")
		}
		return body, nil
	}
}
