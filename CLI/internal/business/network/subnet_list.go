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

func newSubnetListCommand() meta.CommandMeta {
	allowedPlaneTypes := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}}

	return meta.CommandMeta{
		Domain:        "network",
		CanonicalPath: []string{"network", "subnet", "list"},
		Use:           "list",
		Description:   "List subnets",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			rt.BindIntFlag(cmd.Flags(), "plane-type", 3, "Plane type")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			planeType := rt.Int("plane-type")
			if planeType == 0 {
				planeType = 3
			}
			return inputs.RequireOneOfString("plane-type", strconv.Itoa(planeType), allowedPlaneTypes)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			planeType := rt.Int("plane-type")
			if planeType == 0 {
				planeType = 3
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/clusters/v1/%s/subnet?planeType=%d", rt.String("storage-svc-id"), planeType),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/clusters/v1/subnet"}},
				ReadOnly: true,
			}, nil
		},
	}
}
