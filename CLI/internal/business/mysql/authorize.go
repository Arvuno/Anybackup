package mysql

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

type mysqlAuthorizeRequest struct {
	InstanceName string          `json:"instanceName"`
	ClientID     string          `json:"clientId"`
	Username     string          `json:"username"`
	Password     string          `json:"password"`
	SystemID     string          `json:"systemId"`
	ResourceID   string          `json:"resourceId"`
	OSUserName   string          `json:"osUserName"`
	Type         json.RawMessage `json:"type"`
}

func newAuthorizeCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "authorize"},
		Use:           "authorize",
		Description:   "Authorize a MySQL instance",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			return validateMySQLAuthorizeBody(rt.Body)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/resource_center/v1/databasealone/mysql/authorize",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}

func validateMySQLAuthorizeBody(body []byte) error {
	if len(body) == 0 {
		return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
	}

	var req mysqlAuthorizeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	if err := requireAuthorizeField("instanceName", req.InstanceName); err != nil {
		return err
	}
	if err := requireAuthorizeField("clientId", req.ClientID); err != nil {
		return err
	}
	if err := requireAuthorizeField("username", req.Username); err != nil {
		return err
	}
	if err := requireAuthorizeField("password", req.Password); err != nil {
		return err
	}
	if err := requireAuthorizeField("systemId", req.SystemID); err != nil {
		return err
	}
	if err := requireAuthorizeField("resourceId", req.ResourceID); err != nil {
		return err
	}
	if err := requireAuthorizeField("osUserName", req.OSUserName); err != nil {
		return err
	}

	if len(req.Type) == 0 || string(req.Type) == "null" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field type")
	}

	return nil
}

func requireAuthorizeField(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field "+name)
	}
	return nil
}
