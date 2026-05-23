package storage

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPoolNodeListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "storage",
		CanonicalPath: []string{"storage", "pool", "node", "list"},
		Use:           "list",
		Description:   "List storage pool nodes",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 100, "Page size")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			return inputs.ValidatePaging(rt.Int("index"), rt.Int("count"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			count := rt.Int("count")
			if count == 0 {
				count = 100
			}
			query := "index=" + strconv.Itoa(rt.Int("index")) + "&count=" + strconv.Itoa(count)
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/storageresmgm/v1/%s/pool/node?%s", rt.String("storage-svc-id"), query),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/storageresmgm/v1/pool/node"}},
				ReadOnly: true,
			}, nil
		},
	}
}
