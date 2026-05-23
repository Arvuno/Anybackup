package storage

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newServiceListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "storage",
		CanonicalPath: []string{"storage", "service", "list"},
		Use:           "list",
		Description:   "List storage services",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags:     func(cmd *cobra.Command, rt *meta.Runtime) {},
		Validate: func(rt *meta.Runtime) error {
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/mstoragesvcmgm/v1/storage-svc?onlyStorage=true",
				ReadOnly: true,
			}, nil
		},
	}
}
