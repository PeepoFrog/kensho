package sekaiendpoint

type protocolVersion struct {
	P2P   string `json:"p2p"`
	Block string `json:"block"`
	App   string `json:"app"`
}

type other struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

type nodeInfo struct {
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
	CatchingUp          bool   `json:"catching_up"`
}

type pubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type validatorInfo struct {
	Address     string `json:"address"`
	PubKey      pubKey `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type result struct {
	NodeInfo      nodeInfo      `json:"node_info"`
	SyncInfo      syncInfo      `json:"sync_info"`
	ValidatorInfo validatorInfo `json:"validator_info"`
}

type Status struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  result `json:"result"`
}
