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
	Key          string `json:"key"`
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

type Result struct {
	Success    bool          `json:"success"`
	ReturnData []interface{} `json:"return_data"`
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
	multiCaller *core.MultiCaller
}

func NewContractBuilder() ContractBuilder {
	contract := &Contract{
		calls:      make([]core.Call, 0),
		methods:    make([]Method, 0),
		rawMethods: make(map[string]string, 0),
	}

	return contract.WithChainConfig(DefaultChainConfigs[Ethereum])
}

func (ct *Contract) WithChainConfig(config ChainConfig) *Contract {
	if config.MultiCallAddress == "" || config.Url == "" {
		panic("Invalid configuration. MultiCallAddress and Url must be set")
	}

	client, err := ethclient.Dial(config.Url)
	if err != nil {
		panic(err)
	}

	return ct.WithClient(client).AtAddress(config.MultiCallAddress).Build()
}

func (ct *Contract) WithClient(ethClient *ethclient.Client) ContractBuilder {
	ct.ethClient = ethClient
	return ct
}

func (ct *Contract) Build() *Contract {
	return ct
}

func (ct *Contract) AtAddress(address string) ContractBuilder {
	caller, err := core.NewMultiCaller(ct.ethClient, common.HexToAddress(address))
	if err != nil {
		panic(err)
	}
	ct.multiCaller = caller
	return ct
}

func (ct *Contract) AddCall(callName string, contractAddress string, method string, args ...interface{}) *Contract {
	callData, err := ct.contractAbi.Pack(method, args...)
	if err != nil {
		panic(err)
	}
	ct.calls = append(ct.calls, core.Call{
		Method:   method,
		Target:   common.HexToAddress(contractAddress),
		Key:      callName,
		CallData: callData,
	})
	return ct
}

func (ct *Contract) AddMethod(signature string) *Contract {
	existCall, ok := ct.rawMethods[strings.ToLower(signature)]
	if ok {
		panic("MultiCaller named " + existCall + " is exist on ABI")
	}
	ct.rawMethods[strings.ToLower(signature)] = signature
	ct.methods = append(ct.methods, parseNewMethod(signature))
	newAbi, err := repackAbi(ct.methods)
	if err != nil {
		panic(err)
	}
	ct.contractAbi = newAbi
	if err != nil {
		panic(err)
	}
	return ct
}

func (ct *Contract) Abi() abi.ABI {
	return ct.contractAbi
}

func (ct *Contract) Call(blockNumber *big.Int) (*big.Int, map[string][]interface{}, error) {
	res := make(map[string][]interface{})
	blockNumber, results, err := ct.multiCaller.StrictlyExecute(ct.calls, blockNumber)
	if err != nil {
		ct.ClearCall()
		return nil, nil, err
	}
	for _, call := range ct.calls {
		res[call.Key], err = ct.contractAbi.Unpack(call.Method, results[call.Key].ReturnData)
		if err != nil {
			ct.ClearCall()
			return nil, nil, err
		}
	}
	ct.ClearCall()
	return blockNumber, res, err
}

func (ct *Contract) FlexibleCall(requireSuccess bool) (map[string]Result, error) {
	res := make(map[string]Result)
	results, err := ct.multiCaller.Execute(ct.calls, requireSuccess)
	if err != nil {
		ct.ClearCall()
		return nil, err
	}
	for _, call := range ct.calls {
		callSuccess := results[call.Key].Status
		if callSuccess {
			data, err := ct.contractAbi.Unpack(call.Method, results[call.Key].ReturnData)
			if err != nil {
				ct.ClearCall()
				return nil, err
			}
			res[call.Key] = Result{
				Success:    results[call.Key].Status,
				ReturnData: data,
			}
		} else {
			res[call.Key] = Result{
				Success:    results[call.Key].Status,
				ReturnData: nil,
			}
		}
	}
	ct.ClearCall()
	return res, err
}

func (ct *Contract) ClearCall() {
	ct.calls = []core.Call{}
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
						Key:          "",
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
				Key:          "",
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
						Key:          "",
						Type:         strings.TrimSpace(inParam),
						InternalType: strings.TrimSpace(inParam),
					})
				}
			}
		}

		returnType := strings.TrimSpace(singleReturnPaths[1])
		newMethod.Outputs = append(newMethod.Outputs, Argument{
			Key:          "",
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
