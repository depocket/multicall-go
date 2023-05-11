package call

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/depocket/multicall-go/core"
	"github.com/depocket/multicall-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Component struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	InternalType string `json:"internalType"`
}

type Argument struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	InternalType string      `json:"internalType"`
	Components   []Component `json:"components,omitempty"`
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

func (ct *Contract) FlexibleCall(ctx context.Context, requireSuccess bool) (map[string]Result, error) {
	res := make(map[string]Result)
	results, err := ct.multiCaller.Execute(ctx, ct.calls, requireSuccess)
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
	signature = utils.CleanSpaces(signature)
	methodPaths := strings.Split(signature, "(")
	if len(methodPaths) <= 1 {
		panic("Function is invalid format!")
	}
	methodName := strings.Replace(methodPaths[0], "function", "", 1)
	newMethod := Method{
		Name:            methodName,
		Inputs:          make([]Argument, 0),
		Outputs:         make([]Argument, 0),
		Type:            "function",
		StateMutability: "view",
	}

	isMultipleReturn := strings.Contains(signature, ")(")
	if isMultipleReturn {
		multipleReturnPaths := strings.SplitN(signature, ")(", 2)
		paramsPaths := strings.SplitN(multipleReturnPaths[0], "(", 2)
		paramsPath := paramsPaths[1]
		newMethod.Inputs = parseArguments(paramsPath, "input")

		outputPath := strings.TrimSuffix(multipleReturnPaths[1], ")")
		newMethod.Outputs = parseArguments(outputPath, "output")
	} else {
		singleReturnPaths := strings.Split(signature, ")")
		paramsPaths := strings.Split(singleReturnPaths[0], "(")
		params := strings.Split(paramsPaths[1], ",")

		if len(params) > 0 {
			for i, inParam := range params {
				if inParam != "" {
					newMethod.Inputs = append(newMethod.Inputs, Argument{
						Name:         fmt.Sprintf("input%d", i),
						Type:         inParam,
						InternalType: inParam,
					})
				}
			}
		}

		returnType := singleReturnPaths[1]
		newMethod.Outputs = append(newMethod.Outputs, Argument{
			Name:         "output",
			Type:         returnType,
			InternalType: returnType,
		})
	}
	return newMethod
}

func repackAbi(methods []Method) (abi.ABI, error) {
	abiString, err := json.Marshal(methods)
	if err != nil {
		return abi.ABI{}, err
	}
	return abi.JSON(strings.NewReader(string(abiString)))
}

func parseArguments(path, nameFormat string) []Argument {
	result := []Argument{}
	if path == "" {
		return result
	}
	arguments := parsePath(path)
	for i, inArgument := range arguments {
		name := fmt.Sprintf(nameFormat+"%d", i)
		argumentType, isTuple := getArgumentType(inArgument)
		argument := Argument{
			Name:         name,
			Type:         argumentType,
			InternalType: argumentType,
			Components:   make([]Component, 0),
		}

		if isTuple {
			s := strings.ReplaceAll(inArgument, "(", "")
			s = strings.ReplaceAll(s, ")", "")
			s = strings.ReplaceAll(s, "[]", "")
			components := strings.Split(s, ",")
			for j, component := range components {
				argument.Components = append(argument.Components,
					Component{
						Name:         name + fmt.Sprintf("component%d", j),
						Type:         component,
						InternalType: component,
					},
				)
			}
		}
		result = append(result, argument)
	}
	return result
}

func getArgumentType(inArgument string) (string, bool) {
	if strings.Contains(inArgument, "(") {
		if strings.Contains(inArgument, "[]") {
			return "tuple[]", true
		}
		return "tuple", true
	}
	return inArgument, false
}

func parsePath(path string) []string {
	arguments := []string{}
	i := 0
	n := len(path)
	for i < n {
		accumulateString := ""
		char := string(path[i])
		if char == "(" {
			for char != ")" {
				accumulateString += char
				i++
				char = string(path[i])
			}
		}
		for i < n {
			char = string(path[i])
			if char == "," {
				break
			}
			accumulateString += char
			i++
		}
		arguments = append(arguments, accumulateString)
		i++
	}
	return arguments
}
