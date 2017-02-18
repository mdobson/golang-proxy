package data

type ProxyData struct {
	TargetURL       string `yaml:"target_url"`
	Revision        string `yaml:"revision"`
	TargetName      string `yaml:"target_name"`
	ProxyName       string `yaml:"name"`
	Basepath        string `yaml:"base_path"`
	Vhost           string `yaml:"vhost"`
	ProxyConfigName string
	ProxyScope      string
}
