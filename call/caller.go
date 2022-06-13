package call

import (
	"math/big"
	"strings"

	"github.com/depocket/multicall-go/core"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CallerImpl interface {
	WithClient(client *ethclient.Client) CallerImpl
	AtAddress(address string) CallerImpl
	AddCall(callName string, contractAddress string, method string, args ...interface{}) CallerImpl
	AddMethod(signature string) CallerImpl
	Call(blockNumber *big.Int) (*big.Int, map[string][]interface{}, error)
}

type CallerOptions struct {
	Chain   Chain
	Address string
	Url     string
}

type Caller struct {
	ABI         abi.ABI
	Calls       []core.Call
	Client      *ethclient.Client
	Methods     []Method
	MultiCaller *core.MultiCaller
	RawMethods  map[string]string
}

func NewCaller(options *CallerOptions) CallerImpl {
	caller := &Caller{
		Calls:      make([]core.Call, 0),
		Methods:    make([]Method, 0),
		RawMethods: make(map[string]string, 0),
	}

	if options == nil {
		return caller
	}

	var chainInfo ChainInfo
	if options.Chain == "" {
		if options.Address == "" || options.Url == "" {
			panic("Invalid options. Chain or (Address and Url) must be set")
		} else {
			chainInfo = ChainInfo{
				MultiCall: options.Address,
				Url:       options.Url,
			}
		}
	} else {
		// Get default options
		chain, ok := Chains[options.Chain]
		if !ok {
			panic("Invalid options. Chain is not supported")
		}
		// Override options
		if options.Url != "" {
			chain.Url = options.Url
		}
		if options.Address != "" {
			chain.MultiCall = options.Address
		}
		chainInfo = chain
	}

	client, err := ethclient.Dial(chainInfo.Url)
	if err != nil {
		panic(err)
	}

	return caller.WithClient(client).AtAddress(chainInfo.MultiCall)
}

func (c *Caller) WithClient(client *ethclient.Client) CallerImpl {
	c.Client = client
	return c
}

func (c *Caller) AtAddress(address string) CallerImpl {
	caller, err := core.NewMultiCaller(c.Client, common.HexToAddress(address))
	if err != nil {
		panic(err)
	}
	c.MultiCaller = caller
	return c
}

func (c *Caller) AddCall(callName string, contractAddress string, method string, args ...interface{}) CallerImpl {
	callData, err := c.ABI.Pack(method, args...)
	if err != nil {
		panic(err)
	}
	c.Calls = append(c.Calls, core.Call{
		Method:   method,
		Target:   common.HexToAddress(contractAddress),
		Name:     callName,
		CallData: callData,
	})
	return c
}

func (c *Caller) AddMethod(signature string) CallerImpl {
	existCall, ok := c.RawMethods[strings.ToLower(signature)]
	if ok {
		panic("Caller named " + existCall + " is exist on ABI")
	}
	c.RawMethods[strings.ToLower(signature)] = signature
	c.Methods = append(c.Methods, parseNewMethod(signature))
	newAbi, err := repackAbi(c.Methods)
	if err != nil {
		panic(err)
	}
	c.ABI = newAbi
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Caller) Call(blockNumber *big.Int) (*big.Int, map[string][]interface{}, error) {
	res := make(map[string][]interface{})
	blockNumber, results, err := c.MultiCaller.Execute(c.Calls, blockNumber)
	for _, call := range c.Calls {
		res[call.Name], _ = c.ABI.Unpack(call.Method, results[call.Name].ReturnData)
	}
	c.ClearCall()
	return blockNumber, res, err
}

func (c *Caller) ClearCall() {
	c.Calls = []core.Call{}
}
