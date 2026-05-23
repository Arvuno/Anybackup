package job

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
)

func newAppTypesCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "app-types"},
		Use:           "app-types",
		Description:   "List job app types",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/job_center/v1/app_types",
				ReadOnly: true,
			}, nil
		},
	}
}
