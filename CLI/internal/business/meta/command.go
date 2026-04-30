package meta

import (
	"context"

	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

type RiskLevel string

const (
	RiskReadOnly RiskLevel = "read-only"
	RiskWrite    RiskLevel = "write"
)

type AuthType string

const AuthAKSK AuthType = "ak/sk"

type CommandMeta struct {
	Domain        string
	CanonicalPath []string
	Use           string
	Description   string
	Risk          RiskLevel
	AuthType      AuthType
	ReadOnly      bool

	BindFlags    func(cmd *cobra.Command, rt *Runtime)
	Validate     func(rt *Runtime) error
	Execute      func(ctx context.Context, rt *Runtime, deps Dependencies) error
	BuildRequest func(rt *Runtime) (console.RequestSpec, error)
}

type Executor interface {
	Execute(ctx context.Context, rt *Runtime, req console.RequestSpec) ([]byte, error)
}

type ExecutorFunc func(ctx context.Context, rt *Runtime, req console.RequestSpec) ([]byte, error)

func (f ExecutorFunc) Execute(ctx context.Context, rt *Runtime, req console.RequestSpec) ([]byte, error) {
	return f(ctx, rt, req)
}
