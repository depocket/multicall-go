package call

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestAddresses = map[Chain]string{
	Arbitrum:  "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
	Aurora:    "0x8BEc47865aDe3B172A928df8f990Bc7f2A3b9f79",
	Avalanche: "0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7",
	Bsc:       "0x7d99eda556388Ad7743A1B658b9C4FC67D7A9d74",
	Ethereum:  "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
	Fantom:    "0xe1146b9AC456fCbB60644c36Fd3F868A9072fc6E",
	Moonbeam:  "0xeFAeeE334F0Fd1712f9a8cc375f427D9Cdd40d73",
	Moonriver: "0xE3F5a90F9cb311505cd691a46596599aA1A0AD7D",
}

func TestCaller_NewCallerOptions(t *testing.T) {
	assert.PanicsWithValue(t, "Invalid options. Chain or (Address and Url) must be set", func() { NewCaller(&CallerOptions{}) })
	assert.PanicsWithValue(t, "Invalid options. Chain or (Address and Url) must be set", func() { NewCaller(&CallerOptions{Address: "Invalid Address"}) })
	assert.PanicsWithValue(t, "Invalid options. Chain or (Address and Url) must be set", func() { NewCaller(&CallerOptions{Url: "Invalid Url"}) })

	assert.Panics(t, func() { NewCaller(&CallerOptions{Address: "Invalid Address", Url: "Invalid Url"}) })

	assert.PanicsWithValue(t, "Invalid options. Chain is not supported", func() { NewCaller(&CallerOptions{Chain: "Invalid Chain"}) })

	assert.NotPanics(t, func() { NewCaller(&CallerOptions{Chain: Bsc}) })
	assert.NotPanics(t, func() {
		NewCaller(&CallerOptions{
			Chain:   Bsc,
			Address: Chains[Aurora].MultiCall,
			Url:     Chains[Aurora].Url,
		})
	})
}

func TestCaller_Call(t *testing.T) {
	for chain, address := range TestAddresses {
		caller := NewCaller(&CallerOptions{Chain: chain}).AddMethod("function totalSupply()(uint256)")
		_, result, err := caller.AddCall("ts", address, "totalSupply").Call(nil)
		if err != nil {
			assert.Failf(t, "Error calling %s contract", string(chain))
		} else {
			assert.Equal(t, result["ts"][0].(*big.Int).Cmp(big.NewInt(0)), 1)
		}
	}
}

func TestCaller_OverrideOptions(t *testing.T) {
	auroraChain := Chains[Aurora]
	bscChain := Chains[Bsc]

	errBscCaller := NewCaller(&CallerOptions{
		Chain:   Bsc,
		Address: auroraChain.MultiCall,
		Url:     auroraChain.Url,
	}).AddMethod("function totalSupply()(uint256)")
	_, res, err := errBscCaller.AddCall("ts", TestAddresses[Bsc], "totalSupply").Call(nil)
	if err != nil {
		assert.Fail(t, "Error override options")
	} else {
		assert.Equal(t, len(res["ts"]), 0)
	}

	bscCaller := NewCaller(&CallerOptions{
		Chain:   Aurora,
		Address: bscChain.MultiCall,
		Url:     bscChain.Url,
	}).AddMethod("function totalSupply()(uint256)")
	_, res, err = bscCaller.AddCall("ts", TestAddresses[Bsc], "totalSupply").Call(nil)
	if err != nil {
		assert.Fail(t, "Error override options")
	} else {
		assert.Equal(t, res["ts"][0].(*big.Int).Cmp(big.NewInt(0)), 1)
	}
}
