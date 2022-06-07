package call

import (
	"fmt"
	"github.com/depocket/multicall-go/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestAbiBuilder(t *testing.T) {
	client, err := ethclient.Dial("https://bsc-dataseed1.ninicoin.io/")
	demoBep20Abi := NewContractBuilder().WithClient(client).
		AtAddress(BinanceChain).
		AddMethod("function totalSupply()(uint256)").
		Build()
	_, result, err := demoBep20Abi.
		AddCall("supply", "0x7d99eda556388Ad7743A1B658b9C4FC67D7A9d74", "totalSupply").
		Call(nil)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, utils.WeiToEther(result["supply"][0].(*big.Int)).String(), big.NewFloat(21000000).String())
}
