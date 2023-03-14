package core

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/himitsuko/ethclient/ethclient"
	"math/big"
	"strings"
)

type Call struct {
	Key      string         `json:"key"`
	Method   string         `json:"method"`
	Target   common.Address `json:"target"`
	CallData []byte         `json:"call_data"`
}

type CallResponse struct {
	Method     string
	Status     bool
	ReturnData []byte `json:"returnData"`
}

func (call Call) GetMultiCall() MultiCall {
	return MultiCall{Target: call.Target, CallData: call.CallData}
}

type MultiCaller struct {
	Client          *ethclient.Client
	Abi             abi.ABI
	ContractAddress common.Address
}

func NewMultiCaller(client *ethclient.Client, contractAddress common.Address) (*MultiCaller, error) {
	mcAbi, err := abi.JSON(strings.NewReader(MultiMetaData.ABI))
	if err != nil {
		return nil, err
	}

	return &MultiCaller{
		Client:          client,
		Abi:             mcAbi,
		ContractAddress: contractAddress,
	}, nil
}

func (caller *MultiCaller) StrictlyExecute(calls []Call, blockNumber *big.Int) (*big.Int, map[string]CallResponse, error) {
	var multiCalls = make([]MultiCall, 0, len(calls))
	for _, call := range calls {
		multiCalls = append(multiCalls, call.GetMultiCall())
	}
	callData, err := caller.Abi.Pack("aggregate", multiCalls)
	if err != nil {
		return nil, nil, err
	}
	resp, err := caller.Client.CallContract(context.Background(), ethereum.CallMsg{To: &caller.ContractAddress, Data: callData}, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	responses, err := caller.Abi.Unpack("aggregate", resp)

	if err != nil {
		return nil, nil, err
	}

	results := make(map[string]CallResponse)
	for i, response := range responses[1].([][]byte) {
		results[calls[i].Key] = CallResponse{
			Method:     calls[i].Method,
			Status:     true,
			ReturnData: response,
		}
	}
	return responses[0].(*big.Int), results, nil
}

func (caller *MultiCaller) Execute(calls []Call, requireSuccess bool) (map[string]CallResponse, error) {
	var multiCalls = make([]MultiCall, 0, len(calls))
	for _, call := range calls {
		multiCalls = append(multiCalls, call.GetMultiCall())
	}
	callData, err := caller.Abi.Pack("tryAggregate", requireSuccess, multiCalls)
	if err != nil {
		return nil, err
	}
	resp, err := caller.Client.CallContract(context.Background(), ethereum.CallMsg{To: &caller.ContractAddress, Data: callData}, nil)
	if err != nil {
		return nil, err
	}

	responses, err := caller.Abi.Unpack("tryAggregate", resp)

	if err != nil {
		return nil, err
	}

	results := make(map[string]CallResponse)
	for i, response := range responses[0].([]struct {
		Success    bool   `json:"success"`
		ReturnData []byte `json:"returnData"`
	}) {
		results[calls[i].Key] = CallResponse{
			Method:     calls[i].Method,
			ReturnData: response.ReturnData,
			Status:     response.Success,
		}
	}
	return results, nil
}
