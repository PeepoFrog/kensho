package sekai

type ABCI_Info struct {
	Jsonrpc     string      `json:"jsonrpc"`
	ID          int         `json:"id"`
	ABCI_result abci_result `json:"result"`
}

type abci_result struct {
	Response abci_ResponseData `json:"response"`
}

type abci_ResponseData struct {
	Data             string `json:"data"`
	Version          string `json:"version"`
	LastBlockHeight  string `json:"last_block_height"`
	LastBlockAppHash string `json:"last_block_app_hash"`
}
