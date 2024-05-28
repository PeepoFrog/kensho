package httph

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	interxendpoint "github.com/KiraCore/kensho/types/interxEndpoint"
	sekaiendpoint "github.com/KiraCore/kensho/types/sekaiEndpoint"
)

func MakeHttpRequest(url, method string) ([]byte, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetInterxStatus(nodeIP string) (*interxendpoint.Status, error) {
	url := fmt.Sprintf("http://%v:11000/api/status", nodeIP)
	b, err := MakeHttpRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	var info *interxendpoint.Status
	err = json.Unmarshal(b, &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func GetSekaiStatus(nodeIP, port string) (*sekaiendpoint.Status, error) {
	url := fmt.Sprintf("http://%v:%v/status", nodeIP, port)
	b, err := MakeHttpRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	var info *sekaiendpoint.Status
	err = json.Unmarshal(b, &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func ValidatePortRange(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false // Not an integer
	}
	if port < 1 || port > 65535 {
		return false // Out of valid port range
	}
	return true
}

func ValidateIP(input string) bool {
	ipCheck := net.ParseIP(input)
	return ipCheck != nil
}
