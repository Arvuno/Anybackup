package version

type Info struct {
	CLIName                 string   `json:"cliName"`
	CLIVersion              string   `json:"cliVersion"`
	SupportedTargetVersions []string `json:"supportedTargetVersions"`
	DefaultTargetVersion    string   `json:"defaultTargetVersion"`
}
