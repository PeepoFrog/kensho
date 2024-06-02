package shidai

type componentInfo struct {
	Version string `json:"version"`
	Infra   bool   `json:"infra"`
}

type Status struct {
	Sekai  componentInfo `json:"sekai"`
	Interx componentInfo `json:"interx"`
	Shidai componentInfo `json:"shidai"`
	Syslog componentInfo `json:"syslog-ng"`
}

type ValidatorStatus string

const (
	Paused   ValidatorStatus = "PAUSED"
	Active   ValidatorStatus = "ACTIVE"
	Inactive ValidatorStatus = "INACTIVE"
	Jailed   ValidatorStatus = "JAILED"
)

type Validator struct {
	Validator        bool            `json:"validator"`
	ClaimSeat        bool            `json:"claim_seat"`
	ValidatorAddress string          `json:"validator_address"`
	ValidatorStatus  ValidatorStatus `json:"validator_status "`
}
