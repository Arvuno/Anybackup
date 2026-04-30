package inputs

import (
	"net/url"
	"os"
	"strings"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

const DefaultTargetVersion = "9.0.9.0"

type RemoteOptions struct {
	TenantID      string
	Endpoint      string
	AK            string
	SK            string
	TargetVersion string
}

func (o *RemoteOptions) Normalize() error {
	o.TenantID = strings.TrimSpace(o.TenantID)
	if o.TenantID == "" {
		o.TenantID = strings.TrimSpace(os.Getenv("FOUNDATION_TENANT_ID"))
	}
	if o.TargetVersion == "" {
		o.TargetVersion = DefaultTargetVersion
	}
	return o.Validate()
}

func (o RemoteOptions) Validate() error {
	if strings.TrimSpace(o.Endpoint) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --endpoint")
	}
	parsed, err := url.Parse(o.Endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --endpoint, expected full base url")
	}
	if strings.TrimSpace(o.AK) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --ak")
	}
	if strings.TrimSpace(o.SK) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --sk")
	}
	if o.TargetVersion != DefaultTargetVersion {
		return clierrors.New(clierrors.CodeUnsupportedVersion, clierrors.ExitCompatibility, "unsupported --target-version")
	}
	return nil
}
