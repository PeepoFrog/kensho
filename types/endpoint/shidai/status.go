package shidai

type ComponentInfo struct {
	Version string `json:"version"`
	Infra   bool   `json:"infra"`
}

type Status struct {
	Sekai  ComponentInfo `json:"sekai"`
	Interx ComponentInfo `json:"interx"`
	Shidai ComponentInfo `json:"shidai"`
	Syslog ComponentInfo `json:"syslog-ng"`
}
