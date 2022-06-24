package call

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/depocket/multicall-go/core"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Argument struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	InternalType string `json:"internalType"`
}

type Method struct {
	Name            string     `json:"name"`
	Inputs          []Argument `json:"inputs"`
	Outputs         []Argument `json:"outputs"`
	Type            string     `json:"type"`
	StateMutability string     `json:"stateMutability"`
}

type ContractBuilder interface {
	WithClient(ethClient *ethclient.Client) ContractBuilder
	AtAddress(contractAddress string) ContractBuilder
	AddMethod(signature string) *Contract
	Abi() abi.ABI
	Build() *Contract
	WithChainConfig(config ChainConfig) *Contract
}

type Contract struct {
	ethClient   *ethclient.Client
	contractAbi abi.ABI
	rawMethods  map[string]string
	methods     []Method
	calls       []core.Call
	multiCaller *core.Caller
}

func NewContractBuilder() ContractBuilder {
	contract := &Contract{
		calls:      make([]core.Call, 0),
		methods:    make([]Method, 0),
		rawMethods: make(map[string]string, 0),
	}

	return contract.WithChainConfig(DefaultChainConfigs[Ethereum])
}

func (c *Contract) WithChainConfig(config ChainConfig) *Contract {
	if config.MultiCallAddress == "" || config.Url == "" {
		panic("Invalid configuration. MultiCallAddress and Url must be set")
	}

	client, err := ethclient.Dial(config.Url)
	if err != nil {
		panic(err)
	}

	return c.WithClient(client).AtAddress(config.MultiCallAddress).Build()
}

func (a *Contract) WithClient(ethClient *ethclient.Client) ContractBuilder {
	a.ethClient = ethClient
	return a
}

func (a *Contract) Build() *Contract {
	return a
}

func (a *Contract) AtAddress(address string) ContractBuilder {
	caller, err := core.NewCaller(a.ethClient, common.HexToAddress(address))
	if err != nil {
		panic(err)
	}
	a.multiCaller = caller
	return a
}

func (a *Contract) AddCall(callName string, contractAddress string, method string, args ...interface{}) *Contract {
	callData, err := a.contractAbi.Pack(method, args...)
	if err != nil {
		panic(err)
	}
	a.calls = append(a.calls, core.Call{
		Method:   method,
		Target:   common.HexToAddress(contractAddress),
		Name:     callName,
		CallData: callData,
	})
	return a
}

func (a *Contract) AddMethod(signature string) *Contract {
	existCall, ok := a.rawMethods[strings.ToLower(signature)]
	if ok {
		panic("Caller named " + existCall + " is exist on ABI")
	}
	a.rawMethods[strings.ToLower(signature)] = signature
	a.methods = append(a.methods, parseNewMethod(signature))
	newAbi, err := repackAbi(a.methods)
	if err != nil {
		panic(err)
	}
	a.contractAbi = newAbi
	if err != nil {
		panic(err)
	}
	return a
}

func (a *Contract) Abi() abi.ABI {
	return a.contractAbi
}

func (a *Contract) Call(blockNumber *big.Int) (*big.Int, map[string][]interface{}, error) {
	res := make(map[string][]interface{})
	blockNumber, results, err := a.multiCaller.Execute(a.calls, blockNumber)
	for _, call := range a.calls {
		res[call.Name], _ = a.contractAbi.Unpack(call.Method, results[call.Name].ReturnData)
	}
	a.ClearCall()
	return blockNumber, res, err
}

func (a *Contract) ClearCall() {
	a.calls = []core.Call{}
}

func parseNewMethod(signature string) Method {
	methodPaths := strings.Split(signature, "(")
	if len(methodPaths) <= 1 {
		panic("Function is invalid format!")
	}
	methodName := strings.Replace(methodPaths[0], "function", "", 1)
	methodName = strings.TrimSpace(methodName)
	newMethod := Method{
		Name:            methodName,
		Inputs:          make([]Argument, 0),
		Outputs:         make([]Argument, 0),
		Type:            "function",
		StateMutability: "view",
	}

	isMultipleReturn := strings.Contains(signature, ")(")
	if isMultipleReturn {
		multipleReturnPaths := strings.Split(signature, ")(")
		multipleReturnPath := multipleReturnPaths[1]
		paramsPaths := strings.Split(multipleReturnPaths[0], "(")
		params := parseParamsPath(paramsPaths[1])
		if len(params) > 0 {
			for _, inParam := range params {
				if inParam != "" {
					newMethod.Inputs = append(newMethod.Inputs, Argument{
						Name:         "",
						Type:         strings.TrimSpace(inParam),
						InternalType: strings.TrimSpace(inParam),
					})
				}
			}
		}

		outputPath := strings.Replace(multipleReturnPath, ")", "", 1)
		outputs := strings.Split(outputPath, ",")

		for _, outParam := range outputs {
			newMethod.Outputs = append(newMethod.Outputs, Argument{
				Name:         "",
				Type:         strings.TrimSpace(outParam),
				InternalType: strings.TrimSpace(outParam),
			})
		}
	} else {
		singleReturnPaths := strings.Split(signature, ")")
		paramsPaths := strings.Split(singleReturnPaths[0], "(")
		params := parseParamsPath(paramsPaths[1])

		if len(params) > 0 {
			for _, inParam := range params {
				if inParam != "" {
					newMethod.Inputs = append(newMethod.Inputs, Argument{
						Name:         "",
						Type:         strings.TrimSpace(inParam),
						InternalType: strings.TrimSpace(inParam),
					})
				}
			}
		}

		returnType := strings.TrimSpace(singleReturnPaths[1])
		newMethod.Outputs = append(newMethod.Outputs, Argument{
			Name:         "",
			Type:         returnType,
			InternalType: returnType,
		})
	}
	return newMethod
}

func parseParamsPath(paramsPath string) []string {
	params := strings.Split(paramsPath, ",")
	return params
}

func repackAbi(methods []Method) (abi.ABI, error) {
	abiString, err := json.Marshal(methods)
	if err != nil {
		return abi.ABI{}, err
	}
	return abi.JSON(strings.NewReader(string(abiString)))
}
