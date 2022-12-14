package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Remote is an interface, which instances are used to communicate with the perun-cardano-wallet server.
type Remote interface {
	// CallSign is the endpoint for signing data with the perun-cardano-wallet.
	CallSign(SigningRequest) (SigningResponse, error)
	// CallVerify is the endpoint for verifying signatures with the perun-cardano-wallet.
	CallVerify(VerificationRequest) (VerificationResponse, error)
	// CallKeyAvailable is the endpoint for verifying that the connected perun-cardano-wallet has the private key to
	// a given Address.
	CallKeyAvailable(request KeyAvailabilityRequest) (KeyAvailabilityResponse, error)
}

// PerunCardanoWallet is a basic implementation Remote implementation that calls perun-cardano-wallet via http.
type PerunCardanoWallet struct {
	serverAddress string
}

func NewPerunCardanoWallet(addr string) *PerunCardanoWallet {
	return &PerunCardanoWallet{serverAddress: addr}
}

// CallSign computes a Signature for the given SigningRequest via the perun-cardano-wallet server.
func (r *PerunCardanoWallet) CallSign(body SigningRequest) (SigningResponse, error) {
	const signEndpoint = "/sign"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return SigningResponse{}, fmt.Errorf("unable to marshal json for signing request body: %w", err)
	}
	jsonResponse, err := r.callEndpoint(jsonBody, signEndpoint)
	if err != nil {
		return SigningResponse{}, fmt.Errorf("failed to call endpoint: %w", err)
	}
	var result SigningResponse
	if err = json.Unmarshal(jsonResponse, &result); err != nil {
		return SigningResponse{}, fmt.Errorf("failed to unmarshal wallet server response for singing: %w", err)
	}
	return result, nil
}

// CallVerify verifies the (message, signature, public key) tuple in the given VerificationRequest via the
// perun-cardano-wallet server.
func (r *PerunCardanoWallet) CallVerify(body VerificationRequest) (VerificationResponse, error) {
	const verifyEndpoint = "/verify"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return false, fmt.Errorf("unable to marshal json for verification request body: %w", err)
	}
	jsonResponse, err := r.callEndpoint(jsonBody, verifyEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to call endpoint: %w", err)
	}
	var result VerificationResponse
	if err = json.Unmarshal(jsonResponse, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal wallet server response for verification: %w", err)
	}
	return result, nil
}

// CallKeyAvailable queries whether the connected perun-cardano-wallet server has the private key for the public key
// given in the KeyAvailabilityRequest.
func (r *PerunCardanoWallet) CallKeyAvailable(body KeyAvailabilityRequest) (KeyAvailabilityResponse, error) {
	const keyAvailableEndpoint = "/keyAvailable"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return false, fmt.Errorf("unable to marshal json for key availaility request body: %w", err)
	}
	jsonResponse, err := r.callEndpoint(jsonBody, keyAvailableEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to call endpoint: %w", err)
	}
	var result KeyAvailabilityResponse
	if err = json.Unmarshal(jsonResponse, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal wallet server response for key availability %w", err)
	}
	return result, nil
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
