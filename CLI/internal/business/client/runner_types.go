package client

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
)

func newRunnerTypesCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "runner-types"},
		Use:           "runner-types",
		Description:   "Get runner types",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BuildRequest: func(_ *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/commons/client/runnerTypes",
				ReadOnly: true,
			}, nil
		},
	}
}
