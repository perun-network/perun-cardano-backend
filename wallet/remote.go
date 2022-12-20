package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const EndpointSignData = "/sign"
const EndpointSignChannelState = "/signChannelState"
const EndpointVerifyDataSignature = "/verify"
const EndpointVerifyChannelStateSignature = "/verifyChannelState"
const EndpointKeyAvailable = "/keyAvailable"

// Remote is an interface, which instances are used to communicate with the perun-cardano-wallet server.
type Remote interface {
	// CallEndpoint calls the given endpoint with the given body, writing the result to the given result.
	CallEndpoint(endpoint string, body interface{}, result interface{}) error
}

// PerunCardanoWallet is a basic implementation Remote implementation that calls perun-cardano-wallet via http.
type PerunCardanoWallet struct {
	serverAddress string
}

// NewPerunCardanoWallet returns a new PerunCardanoWallet with the given server address.
func NewPerunCardanoWallet(addr string) *PerunCardanoWallet {
	return &PerunCardanoWallet{serverAddress: addr}
}

// CallEndpoint calls the given endpoint on the remote wallet and decodes the json response into the given result.
// `result` must be a pointer.
func (r *PerunCardanoWallet) CallEndpoint(endpoint string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("unable to marshal json body: %w", err)
	}
	jsonResponse, err := r.callEndpoint(jsonBody, endpoint)
	if err != nil {
		return fmt.Errorf("failed to call endpoint: %w", err)
	}
	if err = json.Unmarshal(jsonResponse, result); err != nil {
		return fmt.Errorf("failed to unmarshal wallet server response: %w", err)
	}
	return nil
}

// callEndpoint issues a request to the given endpoint with the given body.
func (r *PerunCardanoWallet) callEndpoint(jsonBody []byte, endpoint string) ([]byte, error) {
	request, err := http.NewRequest("POST", r.serverAddress+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to prepare http request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send http request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to interact with wallet server: %s", response.Status)
		}
		return nil, fmt.Errorf(
			"failed to interact with wallet server: %s with error: %s",
			response.Status,
			string(errorBody),
		)
	}
	jsonResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read server response: %w", err)
	}
	return jsonResponse, nil
}

var _ Remote = &PerunCardanoWallet{}
