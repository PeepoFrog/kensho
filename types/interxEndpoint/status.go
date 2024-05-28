package interxendpoint

type pubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type node struct {
	NodeType        string `json:"node_type"`
	SentryNodeID    string `json:"sentry_node_id"`
	SnapshotNodeID  string `json:"snapshot_node_id"`
	ValidatorNodeID string `json:"validator_node_id"`
	SeedNodeID      string `json:"seed_node_id"`
}

type protocolVersion struct {
	P2P   string `json:"p2p"`
	Block string `json:"block"`
	App   string `json:"app"`
}

type other struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

type statusNodeInfo struct {
	ProtocolVersion protocolVersion `json:"protocol_version"`
	ID              string          `json:"id"`
	ListenAddr      string          `json:"listen_addr"`
	Network         string          `json:"network"`
	Version         string          `json:"version"`
	Channels        string          `json:"channels"`
	Moniker         string          `json:"moniker"`
	Other           other           `json:"other"`
}

type syncInfo struct {
	LatestBlockHash     string `json:"latest_block_hash"`
	LatestAppHash       string `json:"latest_app_hash"`
	LatestBlockHeight   string `json:"latest_block_height"`
	LatestBlockTime     string `json:"latest_block_time"`
	EarliestBlockHash   string `json:"earliest_block_hash"`
	EarliestAppHash     string `json:"earliest_app_hash"`
	EarliestBlockHeight string `json:"earliest_block_height"`
	EarliestBlockTime   string `json:"earliest_block_time"`
}

type validatorInfo struct {
	Address     string `json:"address"`
	PubKey      pubKey `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type interxInfo struct {
	PubKey            pubKey `json:"pub_key"`
	Moniker           string `json:"moniker"`
	KiraAddr          string `json:"kira_addr"`
	KiraPubKey        string `json:"kira_pub_key"`
	FaucetAddr        string `json:"faucet_addr"`
	GenesisChecksum   string `json:"genesis_checksum"`
	ChainID           string `json:"chain_id"`
	Version           string `json:"version"`
	LatestBlockHeight string `json:"latest_block_height"`
	CatchingUp        bool   `json:"catching_up"`
	Node              node   `json:"node"`
}
type Status struct {
	ID            string         `json:"id"`
	InterxInfo    interxInfo     `json:"interx_info"`
	NodeInfo      statusNodeInfo `json:"node_info"`
	SyncInfo      syncInfo       `json:"sync_info"`
	ValidatorInfo validatorInfo  `json:"validator_info"`
}
