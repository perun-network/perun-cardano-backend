package channel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"perun.network/go-perun/channel"
	gpwallet "perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wire"
)

const (
	ContractEndpoint    = "/api/contract"
	ActivateEndpoint    = ContractEndpoint + "/activate"
	InstanceEndpoint    = ContractEndpoint + "/instance"
	WebSocketEndpoint   = "/ws"
	StartEndpointFormat = InstanceEndpoint + "/%s/endpoint/start"
	FundEndpointFormat  = InstanceEndpoint + "/%s/endpoint/fund"
	CloseEndpointFormat = InstanceEndpoint + "/%s/endpoint/close"
)

type PAB struct {
	tokenMap            map[channel.ID]types.ChannelToken
	contractInstanceID  string
	acc                 wallet.RemoteAccount
	subscriptionUrlBase *url.URL
	pabRemote
}

type pabRemote struct {
	pabUrl *url.URL
}

// NewPAB creates a new PAB instance. It expects a host string in the format "host:port" (e.g. "localhost:9080").
func NewPAB(host string, acc wallet.RemoteAccount) (*PAB, error) {
	pabUrl, err := url.Parse("http://" + host)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pab url: %w", err)
	}
	subscriptionUrl, err := url.Parse("ws://" + host + WebSocketEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse subscription url: %w", err)
	}
	return &PAB{
		tokenMap:            make(map[channel.ID]types.ChannelToken),
		acc:                 acc,
		subscriptionUrlBase: subscriptionUrl,
		pabRemote: pabRemote{
			pabUrl: pabUrl,
		},
	}, nil
}

// CallEndpoint calls the given endpoint on the remote wallet and decodes the json response into the given result.
// `result` must be a pointer.
func (r *pabRemote) CallEndpoint(endpoint string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("unable to marshal json body: %w", err)
	}
	jsonResponse, err := r.callEndpoint(jsonBody, endpoint)
	if err != nil {
		return fmt.Errorf("failed to call endpoint: %w", err)
	}
	if result == nil {
		return nil
	}
	if err = json.Unmarshal(jsonResponse, result); err != nil {
		return fmt.Errorf("failed to unmarshal PAB server response: %w", err)
	}
	return nil
}

// callEndpoint issues a request to the given endpoint with the given body.
func (r *pabRemote) callEndpoint(jsonBody []byte, endpoint string) ([]byte, error) {
	request, err := http.NewRequest("POST", r.pabUrl.String()+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to prepare http request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send http request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to interact with PAB server: %s", response.Status)
		}
		return nil, fmt.Errorf(
			"failed to interact with PAB server: %s with error: %s",
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

func (p *PAB) SetChannelToken(id channel.ID, token types.ChannelToken) error {
	emptyToken := types.ChannelToken{}
	if p.tokenMap[id] != emptyToken {
		return errors.New("channel token already set in pab")
	}
	p.tokenMap[id] = token
	return nil
}

func (p *PAB) GetChannelToken(id channel.ID) (types.ChannelToken, error) {
	token := p.tokenMap[id]
	emptyToken := types.ChannelToken{}
	if token == emptyToken {
		return emptyToken, errors.New("channel token not set in pab")
	} else {
		return token, nil
	}
}

func (p *PAB) activateContract() error {
	if p.contractInstanceID != "" {
		return nil
	}
	request := wire.MakePerunActivationBody(p.acc.GetCardanoWalletID())
	var response wire.ContractInstanceID
	err := p.pabRemote.CallEndpoint(ActivateEndpoint, request, &response)
	if err != nil {
		return fmt.Errorf("failed to activate contract: %w", err)
	}
	p.contractInstanceID = response.Decode()
	return nil
}

func (p *PAB) createSubscription(id channel.ID, isPerunSub bool) (*AdjudicatorSub, error) {
	request := wire.MakeAdjudicatorSubscriptionActivationBody(id, p.acc.GetCardanoWalletID())
	var response wire.ContractInstanceID
	err := p.pabRemote.CallEndpoint(ActivateEndpoint, request, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to activate subscription contract: %w", err)
	}
	subUrl := p.subscriptionUrlBase.JoinPath(response.Decode())
	return newAdjudicatorSub(subUrl, id, isPerunSub)
}

// NewInternalSubscription creates a new adjudicator subscription for the given channel. The subscription will return
// internal events. These are more specific to the Cardano implementation of the Perun contract and contain more
// information. The internal events can not be used for anything go-perun related.
// For this use NewPerunEventSubscription instead.
// Care: Only ever interact with the subscription's Next and Err methods from within the same goroutine.
func (p *PAB) NewInternalSubscription(id channel.ID) (*AdjudicatorSub, error) {
	return p.createSubscription(id, false)
}

// NewPerunEventSubscription creates a new adjudicator subscription for the given channel. The subscription will return
// perun events (generealized events compatible with the go-perun core).
// Care: Only ever interact with the subscription's Next and Err methods from within the same goroutine.
func (p *PAB) NewPerunEventSubscription(id channel.ID) (*AdjudicatorSub, error) {
	return p.createSubscription(id, true)
}

// Start issues a request to the PAB to start the channel with the given parameters and initial state.
func (p *PAB) Start(cid channel.ID, params types.ChannelParameters, state types.ChannelState) error {
	if p.contractInstanceID == "" {
		if err := p.activateContract(); err != nil {
			return fmt.Errorf("failed to activate contract: %w", err)
		}
	}
	request := wire.MakeOpenParams(cid, params, state)
	err := p.pabRemote.CallEndpoint(fmt.Sprintf(StartEndpointFormat, p.contractInstanceID), request, nil)
	if err != nil {
		return fmt.Errorf("failed to call endpoint start: %w", err)
	}
	return nil
}

// Fund issues a request to the PAB to fund the channel with the given parameters.
func (p *PAB) Fund(cid channel.ID, index channel.Index) error {
	if p.contractInstanceID == "" {
		if err := p.activateContract(); err != nil {
			return fmt.Errorf("failed to activate contract: %w", err)
		}
	}
	ct, err := p.GetChannelToken(cid)
	if err != nil {
		return fmt.Errorf("failed to fund channel: %w", err)
	}
	request := wire.MakeFundParams(cid, ct, uint16(index))
	err = p.pabRemote.CallEndpoint(fmt.Sprintf(FundEndpointFormat, p.contractInstanceID), request, nil)
	if err != nil {
		return fmt.Errorf("failed to call endpoint fund: %w", err)
	}
	return nil
}

// Abort issues a request to the PAB to abort the channel. This only works on channels that
// are not completely funded yet.
func (p *PAB) Abort() {
	//TODO
}

// Dispute issues a request to the PAB to dispute the channel.
func (p *PAB) Dispute() {
	//TODO
}

// Close issues a request to the PAB to close the channel with the given parameters and final state.
func (p *PAB) Close(id channel.ID, params types.ChannelParameters, state types.ChannelState, sigs []gpwallet.Sig) error {
	if p.contractInstanceID == "" {
		if err := p.activateContract(); err != nil {
			return fmt.Errorf("failed to activate contract: %w", err)
		}
	}
	ct, err := p.GetChannelToken(id)
	if err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	request := wire.MakeCloseParams(id, ct, params, state, sigs)
	err = p.pabRemote.CallEndpoint(fmt.Sprintf(CloseEndpointFormat, p.contractInstanceID), request, nil)
	if err != nil {
		return fmt.Errorf("failed to call endpoint close: %w", err)
	}
	return nil
}

// ForceClose issues a request to the PAB to force close the channel with the given parameters and state. This settles
// the current on-chain state of the channel. One can only force close a channel, if it was disputed beforehand and the
// relative time-lock has expired.
func (p *PAB) ForceClose() {
	//TODO
}
