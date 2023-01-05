package wire_test

import (
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/channel/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestMakeChannelState(t *testing.T) {
	rng := pkgtest.Prng(t)
	testChannelState := test.MakeRandomChannelState(rng)

	expected := wire.ChannelState{
		Balances:  testChannelState.Balances,
		ChannelID: testChannelState.ID,
		Final:     testChannelState.Final,
		Version:   testChannelState.Version,
	}
	actual := wire.MakeChannelState(testChannelState)
	require.Equal(t, expected, actual, "wire.ChannelState not as expected")
}

func TestChannelState_Decode(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomChannelState(rng)
	uut := wire.ChannelState{
		Balances:  expected.Balances,
		ChannelID: expected.ID,
		Final:     expected.Final,
		Version:   expected.Version,
	}
	actual := uut.Decode()
	require.Equal(t, expected, actual, "decoded channel state is wrong")
}
func TestChannelState(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomChannelState(rng)
	actual := wire.MakeChannelState(expected).Decode()
	require.Equal(t, expected, actual, "channel state not as expected")
}
