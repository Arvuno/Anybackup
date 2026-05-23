package mysql

import (
	"fmt"
	"strconv"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

func validateCountRange(flag string, value, min, max int) error {
	if value < min || value > max {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateOptionalNonNegativeInt64(flag, value string) error {
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}
