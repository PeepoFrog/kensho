package types

const (
	BOOTSTRAP_SCRIPT       string = "https://raw.githubusercontent.com/KiraCore/sekin/main/scripts/bootstrap.sh"
	SEKIN_EXECUTE_ENDPOINT string = "http://localhost:8282/api/execute"
	SEKIN_STATUS_ENDPOINT  string = "http://localhost:8282/api/execute"

	DEFAULT_INTERX_PORT int = 11000
	DEFAULT_P2P_PORT    int = 26656
	DEFAULT_RPC_PORT    int = 26657
	DEFAULT_GRPC_PORT   int = 9090
	DEFAULT_SHIDAI_PORT int = 8282
)

type RequestDeployPayload struct {
	Command string `json:"command"`
	Args    Args   `json:"args"`
}

// Args represents the arguments in the JSON payload.
type Args struct {
	IP         string `json:"ip"`
	InterxPort string `json:"interx_port"`
	RPCPort    string `json:"rpc_port"`
	P2PPort    string `json:"p2p_port"`
	Mnemonic   string `json:"mnemonic"`
	Local      bool   `json:"local"`
}
