package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	gpwallet "perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
	"polycry.pt/poly-go/sync"
)

type validSignatures struct {
	DataLock                    sync.Mutex
	ChannelStateLock            sync.Mutex
	validDataSignatures         []DataSignature
	validChannelStateSignatures []ChannelStateSignature
}

var ValidSignatures = &validSignatures{
	validDataSignatures:         []DataSignature{},
	validChannelStateSignatures: []ChannelStateSignature{},
}

func AddDataSignature(address address.Address, sig gpwallet.Sig, msg []byte) {
	ValidSignatures.DataLock.Lock()
	defer ValidSignatures.DataLock.Unlock()
	ValidSignatures.validDataSignatures = append(ValidSignatures.validDataSignatures, DataSignature{
		Address:   address,
		Signature: sig,
		Message:   msg,
	})
}

func AddChannelStateSignature(address address.Address, sig gpwallet.Sig, channelState types.ChannelState) {
	ValidSignatures.ChannelStateLock.Lock()
	defer ValidSignatures.ChannelStateLock.Unlock()
	ValidSignatures.validChannelStateSignatures = append(ValidSignatures.validChannelStateSignatures, ChannelStateSignature{
		Address:      address,
		Signature:    sig,
		ChannelState: channelState,
	})
}

func VerifyDataSig(sig DataSignature) bool {
	ValidSignatures.DataLock.Lock()
	defer ValidSignatures.DataLock.Unlock()
	for _, valid := range ValidSignatures.validDataSignatures {
		if sig.Equal(valid) {
			return true
		}
	}
	return false
}

func VerifyChannelStateSig(sig ChannelStateSignature) bool {
	ValidSignatures.ChannelStateLock.Lock()
	defer ValidSignatures.ChannelStateLock.Unlock()
	for _, valid := range ValidSignatures.validChannelStateSignatures {
		if sig.Equal(valid) {
			return true
		}
	}
	return false
}

type DataSignature struct {
	Address   address.Address
	Signature gpwallet.Sig
	Message   []byte
}

func (d DataSignature) Equal(other DataSignature) bool {
	return d.Address.Equal(&other.Address) &&
		bytes.Equal(d.Signature, other.Signature) &&
		bytes.Equal(d.Message, other.Message)
}

type ChannelStateSignature struct {
	Address      address.Address
	Signature    gpwallet.Sig
	ChannelState types.ChannelState
}

func (d ChannelStateSignature) Equal(other ChannelStateSignature) bool {
	return d.Address.Equal(&other.Address) &&
		bytes.Equal(d.Signature, other.Signature) &&
		d.ChannelState.Equal(other.ChannelState)
}

// GenericRemote is a mock remote that suits the generic go-perun wallet tests.
// Signatures are generated randomly and all valid signature t-tuples are collected globally
// in ValidSignatures. This means that GenericRemote verifies a signature as valid, iff it has been signed by ANY
// Generic remote before.
// Note:
// Generic remote does not implement the ChannelState-specific endpoints wallet.EndpointSignChannelState and
// wallet.EndpointVerifyChannelStateSignature because the generic go-perun wallet tests do not test for their use.
type GenericRemote struct {
	AvailableAddress address.Address
	rng              *rand.Rand
	// The lock is needed because the same GenericRemote instance might be required to generate random signatures
	// using the same rand.Rand instance in parallel and rand.Rand is not concurrently safe.
	mutex sync.Mutex

	callEndpoint func(string, interface{}, interface{}) error
}

// NewGenericRemote returns a new generic remote instance. Every GenericRemote instance needs to receive on an exclusive
// rand.Rand instance for concurrent safety!
func NewGenericRemote(availableAddress address.Address, rng *rand.Rand) *GenericRemote {
	g := GenericRemote{
		AvailableAddress: availableAddress,
		rng:              rng,
	}
	g.callEndpoint = makeGenericCallEndpointDefault(&g)
	return &g
}

func (g *GenericRemote) CallEndpoint(endpoint string, request interface{}, response interface{}) error {
	return g.callEndpoint(endpoint, request, response)
}

func (g *GenericRemote) SetCallEndpoint(f func(string, interface{}, interface{}) error) {
	g.callEndpoint = f
}

func makeGenericCallEndpointDefault(g *GenericRemote) func(string, interface{}, interface{}) error {
	return func(endpoint string, req interface{}, resp interface{}) error {
		switch endpoint {
		case wallet.EndpointSignData:
			request, ok := req.(wire.SigningRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to SingingRequest")
			}
			response, ok := resp.(*wire.SigningResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to SingingResponse")
			}
			reqAddr, err := request.PubKey.Decode()
			if err != nil {
				return fmt.Errorf("unable to decode PubKey from request")
			}
			if !reqAddr.Equal(&g.AvailableAddress) {
				return fmt.Errorf("account is not available in wallet")
			}
			msg, err := hex.DecodeString(request.Message)
			if err != nil {
				return fmt.Errorf("unable to decode message")
			}

			g.mutex.Lock()
			sig := MakeRandomSignature(g.rng)
			g.mutex.Unlock()
			AddDataSignature(reqAddr, sig, msg)
			*response = wire.MakeSignature(sig)
			return nil
		case wallet.EndpointVerifyDataSignature:
			request, ok := req.(wire.VerificationRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to VerificationRequest")
			}
			response, ok := resp.(*wire.VerificationResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to VerificationResponse")
			}
			reqAddr, err := request.PubKey.Decode()
			if err != nil {
				return fmt.Errorf("unable to decode PubKey from request")
			}
			sig, err := request.Signature.Decode()
			if err != nil {
				return fmt.Errorf("unable to decode signature")
			}
			msg, err := hex.DecodeString(request.Message)
			if err != nil {
				return fmt.Errorf("unable to decode message")
			}
			*response = VerifyDataSig(DataSignature{
				Address:   reqAddr,
				Signature: sig,
				Message:   msg,
			})
			return nil
		case wallet.EndpointKeyAvailable:
			request, ok := req.(wire.KeyAvailabilityRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to KeyAvailabilityRequst")
			}
			response, ok := resp.(*wire.KeyAvailabilityResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to KeyAvailabilityResponse")
			}
			reqAddr, err := request.Decode()
			if err != nil {
				return fmt.Errorf("unable to decode PubKey from request")
			}
			*response = reqAddr.Equal(&g.AvailableAddress)
			return nil
		default:
			panic("unimplemented endpoint: " + endpoint)
		}
	}
}
