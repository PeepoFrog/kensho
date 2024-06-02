package httph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	interxendpoint "github.com/KiraCore/kensho/types/endpoint/interx"
	sekaiendpoint "github.com/KiraCore/kensho/types/endpoint/sekai"
	shidaiendpoint "github.com/KiraCore/kensho/types/endpoint/shidai"
	"golang.org/x/crypto/ssh"
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

func GetShidaiStatus(sshClient *ssh.Client, shidaiPort int) (*shidaiendpoint.Status, error) {
	valid := ValidatePortRange(strconv.Itoa(shidaiPort))
	if !valid {
		return nil, fmt.Errorf("<%v> is not valid", shidaiPort)
	}
	o, err := ExecHttpRequestBySSHTunnel(sshClient, fmt.Sprintf("http://localhost:%v/status", shidaiPort), "GET", nil)
	if err != nil {
		return nil, err
	}
	var data shidaiendpoint.Status
	err = json.Unmarshal(o, &data)
	if err != nil {
		return nil, err
	}
	return &data, err
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

func ExecHttpRequestBySSHTunnel(sshClient *ssh.Client, address, method string, payload []byte) ([]byte, error) {
	dialer := func(network, addr string) (net.Conn, error) {
		conn, err := sshClient.Dial(network, addr)
		if err != nil {
			log.Printf("Failed to establish SSH tunnel: %v", err)
		}
		return conn, err
	}

	httpTransport := &http.Transport{
		Dial: dialer,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	var req *http.Request
	var err error

	if len(payload) == 0 {
		req, err = http.NewRequest(method, address, nil)
	} else {
		req, err = http.NewRequest(method, address, bytes.NewBuffer(payload))
	}
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return out, nil
}
