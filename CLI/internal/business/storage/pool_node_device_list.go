package storage

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPoolNodeDeviceListCommand() meta.CommandMeta {
	defaultTypes := []string{"1", "2", "3", "4", "5", "7"}
	allowedTypes := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "7": {}}

	return meta.CommandMeta{
		Domain:        "storage",
		CanonicalPath: []string{"storage", "pool", "node", "device", "list"},
		Use:           "list",
		Description:   "List storage devices for one node",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "node-id", "", "Node ID")
			rt.BindStringSliceFlag(cmd.Flags(), "types", append([]string(nil), defaultTypes...), "Device type filter (repeatable)")
			rt.BindIntFlag(cmd.Flags(), "authorized", 2, "Authorization filter")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			if values := rt.Strings("types"); len(values) > 0 {
				if err := inputs.RequireAllInSet("types", values, allowedTypes); err != nil {
					return err
				}
			}
			return validateNodeID(rt)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			types := rt.Strings("types")
			if len(types) == 0 {
				types = defaultTypes
			}
			authorized := rt.Int("authorized")
			if authorized == 0 {
				authorized = 2
			}
			var query strings.Builder
			query.WriteString("nodeId=")
			query.WriteString(rt.String("node-id"))
			for _, value := range types {
				query.WriteString("&types=")
				query.WriteString(value)
			}
			query.WriteString("&authorized=")
			query.WriteString(fmt.Sprintf("%d", authorized))
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/storageresmgm/v1/%s/device?%s", rt.String("storage-svc-id"), query.String()),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/storageresmgm/v1/device"}},
				ReadOnly: true,
			}, nil
		},
	}
}
