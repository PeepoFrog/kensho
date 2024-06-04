package interx

import "time"

type Validators struct {
	Validators []Validator `json:"validators"`
}

type Validator struct {
	Top                   string    `json:"top"`
	Address               string    `json:"address"`
	Valkey                string    `json: "valkey"`
	Pubkey                string    `json:"pubkey"`
	Proposer              string    `json:"proposer"`
	Moniker               string    `json:"moniker"`
	Status                string    `json:"status"`
	Rank                  int       `json:"rank,string"`
	Streak                int       `json:"streak,string"`
	Mischance             int       `json:"mischance,string"`
	MischanceConfidence   int       `json:"mischance_confidence,string"`
	StartHeight           int       `json:"start_height,string"`
	InactiveUntil         time.Time `json:"inactive_until"`
	LastPresentBlock      int       `json:"last_present_block,string"`
	MissedBlocksCounter   int       `json:"missed_blocks_counter,string"`
	ProducedBlocksCounter int       `json:"produced_blocks_counter,string"`
	StakingPoolID         string    `json:"staking_pool_id"`
	StakingPoolStatus     string    `json:"staking_pool_status"`
	Description           string    `json:"description,omitempty"`
	Website               string    `json:"website,omitempty"`
	Logo                  string    `json:"logo,omitempty"`
	Social                string    `json:"social,omitempty"`
	Contact               string    `json:"contact,omitempty"`
}
