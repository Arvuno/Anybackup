package vmware

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newTimepointMetadataCommand() meta.CommandMeta {
	allowedBusiness := map[string]struct{}{"1": {}}

	return meta.CommandMeta{
		Domain:        "vmware",
		CanonicalPath: []string{"vmware", "timepoint", "metadata"},
		Use:           "metadata",
		Description:   "Get VMware timepoint metadata",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "timestamp", "", "Timepoint timestamp")
			rt.BindStringFlag(cmd.Flags(), "data-set-id", "", "Dataset ID")
			rt.BindStringFlag(cmd.Flags(), "business", "", "Business type")
			rt.BindStringFlag(cmd.Flags(), "request-id", "", "Request ID")
			rt.BindStringFlag(cmd.Flags(), "storage-service-id", "", "Storage service ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("timestamp", rt.String("timestamp")); err != nil {
				return err
			}
			if err := inputs.RequireNonEmpty("data-set-id", rt.String("data-set-id")); err != nil {
				return err
			}
			if v := rt.String("business"); v != "" {
				if err := inputs.RequireOneOfString("business", v, allowedBusiness); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder().
				Add("timestamp", rt.String("timestamp")).
				Add("dataSetId", rt.String("data-set-id"))
			if v := rt.String("business"); v != "" {
				qb.Add("business", v)
			}
			if v := rt.String("request-id"); v != "" {
				qb.Add("requestId", v)
			}
			if v := rt.String("storage-service-id"); v != "" {
				qb.Add("storageServiceId", v)
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/virtual/vmware/time_point/get_metadata?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
