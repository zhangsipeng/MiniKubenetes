package runtimedata

type YamlConfig struct {
	APIServerIP   string
	ClientIP      string
	RootCAPem     string
	ClientCertPem string
	ClientKeyPem  string
	Others        map[string]string
}

type RuntimeConfig struct {
	YamlConfig *YamlConfig
	// TODO
}
