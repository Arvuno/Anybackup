package network

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newSubnetNodeListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "network",
		CanonicalPath: []string{"network", "subnet", "node", "list"},
		Use:           "list",
		Description:   "List subnet nodes",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "subnet-id", "", "Subnet ID")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
			rt.BindIntFlag(cmd.Flags(), "plane-type", 3, "Plane type")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			if err := validateSubnetID(rt); err != nil {
				return err
			}
			return inputs.ValidatePaging(rt.Int("index"), rt.Int("count"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			count := rt.Int("count")
			if count == 0 {
				count = 30
			}
			planeType := rt.Int("plane-type")
			if planeType == 0 {
				planeType = 3
			}
			qb := inputs.NewQueryBuilder().
				Add("index", strconv.Itoa(rt.Int("index"))).
				Add("count", strconv.Itoa(count)).
				Add("planeType", strconv.Itoa(planeType)).
				Add("subnetId", rt.String("subnet-id"))
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/clusters/v1/%s/subnet/nodes?%s", rt.String("storage-svc-id"), qb.Encode()),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/clusters/v1/subnet/nodes"}},
				ReadOnly: true,
			}, nil
		},
	}
}
