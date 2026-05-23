package cmd

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/api"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/client"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/host"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/job"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/mysql"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/network"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/policy"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/protect"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/storage"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/timepoint"
	bizversion "github.com/anybackup-ai/Anybackup/CLI/internal/business/version"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/vmware"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/signer"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/transport"

	"github.com/spf13/cobra"
)

type Dependencies struct {
	Streams output.Streams
	Version bizversion.Info
	Console meta.Executor
}

func BuildVersionInfo(cliVersion string) bizversion.Info {
	return bizversion.Info{
		CLIName:                 "foundation-cli",
		CLIVersion:              cliVersion,
		SupportedTargetVersions: []string{"9.0.9.0"},
		DefaultTargetVersion:    "9.0.9.0",
	}
}

func NewRootCommand(deps Dependencies) *cobra.Command {
	root := &cobra.Command{
		Use:           "foundation-cli",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(deps.Streams.Stdout)
	root.SetErr(deps.Streams.Stderr)

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print foundation-cli version metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.WriteJSON(output.Streams{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}, deps.Version)
		},
	})

	defs := append([]meta.CommandMeta{}, api.Commands()...)
	defs = append(defs, policy.Commands()...)
	defs = append(defs, protect.Commands()...)
	defs = append(defs, timepoint.Commands()...)
	defs = append(defs, job.Commands()...)
	defs = append(defs, host.Commands()...)
	defs = append(defs, client.Commands()...)
	defs = append(defs, vmware.Commands()...)
	defs = append(defs, mysql.Commands()...)
	defs = append(defs, network.Commands()...)
	defs = append(defs, storage.Commands()...)

	meta.Register(root, meta.Dependencies{
		Streams: deps.Streams,
		Console: deps.Console,
	}, defs...)

	return root
}

func NewDefaultExecutor(client *http.Client) meta.Executor {
	httpClient := withSDKTLS(client)
	return meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
		exec := console.Executor{
			Signer: signer.Signer{},
			Client: transport.Client{Doer: httpClient},
		}
		return exec.Execute(ctx, rt.Remote, req)
	})
}

func withSDKTLS(client *http.Client) *http.Client {
	if client == nil {
		client = &http.Client{}
	}

	clone := *client
	baseTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &clone
	}
	transport := baseTransport.Clone()
	if existing, ok := client.Transport.(*http.Transport); ok {
		transport = existing.Clone()
	} else if client.Transport != nil {
		return &clone
	}

	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	} else {
		transport.TLSClientConfig = transport.TLSClientConfig.Clone()
	}
	transport.TLSClientConfig.InsecureSkipVerify = true
	clone.Transport = transport
	return &clone
}
