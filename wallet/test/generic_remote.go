// Copyright 2022, 2023 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
type GenericRemote struct {
	AvailableAddress address.Address
	rng              *rand.Rand
	// The lock is needed because the same GenericRemote instance might be required to generate random signatures
	// using the same rand.Rand instance in parallel and rand.Rand is not concurrently safe.
	mutex sync.Mutex

	callEndpoint func(string, interface{}, interface{}) error
}

// NewGenericRemote returns a new generic remote instance. Every GenericRemote instance needs to receive an exclusive
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
				return fmt.Errorf("unable to cast response to SingingResponse")
			}
			return g.endpointSignData(request, response)
		case wallet.EndpointVerifyDataSignature:
			request, ok := req.(wire.VerificationRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to VerificationRequest")
			}
			response, ok := resp.(*wire.VerificationResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to VerificationResponse")
			}
			return g.endpointVerifyDataSignature(request, response)
		case wallet.EndpointKeyAvailable:
			request, ok := req.(wire.KeyAvailabilityRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to KeyAvailabilityRequst")
			}
			response, ok := resp.(*wire.KeyAvailabilityResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to KeyAvailabilityResponse")
			}
			return g.endpointKeyAvailable(request, response)
		case wallet.EndpointSignChannelState:
			request, ok := req.(wire.ChannelStateSigningRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to ChannelStateSigningRequest")
			}
			response, ok := resp.(*wire.SigningResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to SigningResponse")
			}
			return g.endpointSignChannelState(request, response)
		case wallet.EndpointVerifyChannelStateSignature:
			request, ok := req.(wire.ChannelStateVerificationRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to ChannelStateVerificationRequest")
			}
			response, ok := resp.(*wire.VerificationResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to VerificationResponse")
			}
			return g.endpointVerifyChannelStateSignature(request, response)
		default:
			return fmt.Errorf("invalid endpoint: %s", endpoint)
		}
	}
}

func (g *GenericRemote) endpointSignData(request wire.SigningRequest, response *wire.SigningResponse) error {
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
}

func (g *GenericRemote) endpointVerifyDataSignature(request wire.VerificationRequest, response *wire.VerificationResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	sig, err := request.Signature.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode signature from request")
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
}

func (g *GenericRemote) endpointKeyAvailable(request wire.KeyAvailabilityRequest, response *wire.KeyAvailabilityResponse) error {
	reqAddr, err := request.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	*response = reqAddr.Equal(&g.AvailableAddress)
	return nil
}

func (g *GenericRemote) endpointSignChannelState(request wire.ChannelStateSigningRequest, response *wire.SigningResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	if !reqAddr.Equal(&g.AvailableAddress) {
		return fmt.Errorf("account is not available in wallet")
	}
	state := request.ChannelState.Decode()
	g.mutex.Lock()
	sig := MakeRandomSignature(g.rng)
	g.mutex.Unlock()
	AddChannelStateSignature(reqAddr, sig, state)
	*response = wire.MakeSignature(sig)
	return nil
}

func (g *GenericRemote) endpointVerifyChannelStateSignature(request wire.ChannelStateVerificationRequest, response *wire.VerificationResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	sig, err := request.Signature.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode signature from request")
	}
	state := request.ChannelState.Decode()
	*response = VerifyChannelStateSig(ChannelStateSignature{
		Address:      reqAddr,
		Signature:    sig,
		ChannelState: state,
	})
	return nil
}
