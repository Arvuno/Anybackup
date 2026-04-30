package job

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
)

func newBusinessTypesCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "business-types"},
		Use:           "business-types",
		Description:   "List job business types",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/job_center/v1/business_types",
				ReadOnly: true,
			}, nil
		},
	}
}
