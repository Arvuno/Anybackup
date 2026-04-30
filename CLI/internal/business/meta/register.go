package meta

import (
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

type Dependencies struct {
	Streams output.Streams
	Console Executor
}

func Register(root *cobra.Command, deps Dependencies, defs ...CommandMeta) {
	for _, def := range defs {
		registerOne(root, deps, def)
	}
}

func registerOne(root *cobra.Command, deps Dependencies, def CommandMeta) {
	rt := NewRuntime()
	cmd := &cobra.Command{
		Use:   def.Use,
		Short: def.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			rt.CaptureChangedFlags(cmd.Flags())

			if err := rt.Remote.Normalize(); err != nil {
				return err
			}

			if err := loadBodyIfNeeded(rt); err != nil {
				return err
			}

			if def.Validate != nil {
				if err := def.Validate(rt); err != nil {
					return err
				}
			}
			if def.Execute != nil {
				return def.Execute(cmd.Context(), rt, deps)
			}
			req, err := def.BuildRequest(rt)
			if err != nil {
				return err
			}
			if deps.Console == nil {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "console executor is not configured")
			}
			body, err := deps.Console.Execute(cmd.Context(), rt, req)
			if err != nil {
				return err
			}
			return output.WriteRaw(deps.Streams, body)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&rt.Remote.TenantID, "tenant-id", "", "tenant id")
	flags.StringVar(&rt.Remote.Endpoint, "endpoint", "", "backend base url")
	flags.StringVar(&rt.Remote.AK, "ak", "", "access key")
	flags.StringVar(&rt.Remote.SK, "sk", "", "secret key")
	flags.StringVar(&rt.Remote.TargetVersion, "target-version", "", "backend version, defaults to 9.0.9.0")

	if def.BindFlags != nil {
		def.BindFlags(cmd, rt)
	}

	parent := root
	path := def.CanonicalPath
	if len(path) == 0 {
		path = []string{def.Use}
	}
	for _, name := range path[:len(path)-1] {
		parent = ensureSubcommand(parent, name)
	}
	parent.AddCommand(cmd)
}

func loadBodyIfNeeded(rt *Runtime) error {
	if rt.Inline == "" && rt.BodyFile == "" {
		return nil
	}
	body, err := inputs.LoadBody(rt.Inline, rt.BodyFile)
	if err != nil {
		return err
	}
	rt.Body = body
	return nil
}

func ensureSubcommand(parent *cobra.Command, use string) *cobra.Command {
	for _, child := range parent.Commands() {
		if child.Name() == use {
			return child
		}
	}
	next := &cobra.Command{
		Use:   use,
		Short: strings.Title(use),
	}
	parent.AddCommand(next)
	return next
}
