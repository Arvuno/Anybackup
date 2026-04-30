package inputs

import (
	"fmt"
	"net"
	"strings"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

func RequireIP(flag, value string) error {
	if strings.TrimSpace(value) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("missing --%s", flag))
	}
	if net.ParseIP(value) == nil {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}
