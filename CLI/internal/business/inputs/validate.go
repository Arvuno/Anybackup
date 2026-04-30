package inputs

import (
	"fmt"
	"strings"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

func RequireNonEmpty(flag, value string) error {
	if strings.TrimSpace(value) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("missing --%s", flag))
	}
	return nil
}

func ValidatePaging(index, count int) error {
	if index < 0 {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "--index must be non-negative")
	}
	if count <= 0 {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "--count must be positive")
	}
	return nil
}

func RequireOneOfString(flag, value string, allow map[string]struct{}) error {
	if _, ok := allow[value]; !ok {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s value %q", flag, value))
	}
	return nil
}

func RequireAllInSet(flag string, values []string, allow map[string]struct{}) error {
	for _, v := range values {
		if _, ok := allow[v]; !ok {
			return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s value %q", flag, v))
		}
	}
	return nil
}
