package api

import (
	"net/http"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newAPICommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "api",
		CanonicalPath: []string{"api"},
		Use:           "api",
		Description:   "Signed passthrough for a relative AB8 path",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "method", http.MethodGet, "http method")
			rt.BindStringFlag(cmd.Flags(), "path", "", "relative backend path")
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			method := strings.ToUpper(rt.String("method"))
			if !isValidMethod(method) {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --method")
			}
			path := rt.String("path")
			if path == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --path")
			}
			if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || !strings.HasPrefix(path, "/") {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "api only accepts relative path")
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			method := strings.ToUpper(rt.String("method"))
			return console.RequestSpec{
				Method:   method,
				Path:     rt.String("path"),
				Body:     rt.Body,
				ReadOnly: method == http.MethodGet,
			}, nil
		},
	}
}

func isValidMethod(method string) bool {
	if method == "" {
		return false
	}
	for _, r := range method {
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		switch r {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}
