package network

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newSubnetNodeIPSetCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "network",
		CanonicalPath: []string{"network", "subnet", "node", "ip", "set"},
		Use:           "set",
		Description:   "Set subnet node business IP",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "subnet-id", "", "Subnet ID")
			rt.BindStringFlag(cmd.Flags(), "node-id", "", "Node ID")
			rt.BindStringFlag(cmd.Flags(), "ip", "", "Business IP address")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			if err := validateSubnetID(rt); err != nil {
				return err
			}
			if err := validateNodeID(rt); err != nil {
				return err
			}
			return inputs.RequireIP("ip", rt.String("ip"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			body, err := buildNodeIPPayload(rt)
			if err != nil {
				return console.RequestSpec{}, err
			}
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     fmt.Sprintf("/clusters/v1/%s/subnet/nodes/node_ip_addresses", rt.String("storage-svc-id")),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/clusters/v1/subnet/nodes/node_ip_addresses"}},
				Body:     body,
				ReadOnly: false,
			}, nil
		},
	}
}
