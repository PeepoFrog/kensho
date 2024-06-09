package types

const (
	BOOTSTRAP_SCRIPT           string = "https://raw.githubusercontent.com/KiraCore/sekin/main/scripts/bootstrap.sh"
	SEKIN_EXECUTE_ENDPOINT     string = "http://localhost:8282/api/execute"
	SEKIN_EXECUTE_CMD_ENDPOINT string = "http://localhost:8282/api/execute/tx"
	SEKIN_STATUS_ENDPOINT      string = "http://localhost:8282/api/status"

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
	InterxPort int    `json:"interx_port"`
	RPCPort    int    `json:"rpc_port"`
	P2PPort    int    `json:"p2p_port"`
	Mnemonic   string `json:"mnemonic"`
	Local      bool   `json:"local"`
}

type Cmd string

const (
	Activate           Cmd = "activate"
	Pause              Cmd = "pause"
	Unpause            Cmd = "unpause"
	ClaimValidatorSeat Cmd = "claim_seat"
)

type RequestTXPayload struct {
	Command string       `json:"command"`
	Args    ExecSekaiCmd `json:"args"`
}
type ExecSekaiCmd struct {
	TX      Cmd    `json:"tx"` //pause, unpause, activate,
	Moniker string `json:"moniker"`
}
