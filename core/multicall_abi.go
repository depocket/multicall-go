// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package core

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MultiCall is an auto generated low-level Go binding around an user-defined struct.
type MultiCall struct {
	Target   common.Address
	CallData []byte
}

// Result is an auto generated low-level Go binding around an user-defined struct.
type Result struct {
	Success    bool
	ReturnData []byte
}

// MultiMetaData contains all meta data concerning the Multi contract.
var MultiMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"blockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryAggregate\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryBlockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structDePocketCore.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Multi is an auto generated Go binding around an Ethereum contract.
type Multi struct {
	MultiCaller     // Read-only binding to the contract
	MultiTransactor // Write-only binding to the contract
	MultiFilterer   // Log filterer for contract events
}

// MultiCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type MultiTransactorSession struct {
	Contract     *MultiTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

type MultiRaw struct {
	Contract *Multi // Generic contract binding to access the raw methods on
}

type MultiCallerRaw struct {
	Contract *MultiCaller // Generic read-only contract binding to access the raw methods on
}

type MultiTransactorRaw struct {
	Contract *MultiTransactor // Generic write-only contract binding to access the raw methods on
}

func NewMulti(address common.Address, backend bind.ContractBackend) (*Multi, error) {
	contract, err := bindMulti(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multi{MultiCaller: MultiCaller{contract: contract}, MultiTransactor: MultiTransactor{contract: contract}, MultiFilterer: MultiFilterer{contract: contract}}, nil
}

func NewMultiCaller(address common.Address, caller bind.ContractCaller) (*MultiCaller, error) {
	contract, err := bindMulti(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiCaller{contract: contract}, nil
}

func bindMulti(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Multi *MultiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multi.Contract.MultiCaller.contract.Call(opts, result, method, params...)
}

func (_Multi *MultiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multi.Contract.contract.Call(opts, result, method, params...)
}

func (_Multi *MultiCaller) Aggregate(opts *bind.CallOpts, calls []MultiCall) (struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}, error) {
	var out []interface{}
	err := _Multi.contract.Call(opts, &out, "aggregate", calls)

	outstruct := new(struct {
		BlockNumber *big.Int
		ReturnData  [][]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ReturnData = *abi.ConvertType(out[1], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

func (_Multi *MultiCaller) BlockAndAggregate(opts *bind.CallOpts, calls []MultiCall) (struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	ReturnData  []Result
}, error) {
	var out []interface{}
	err := _Multi.contract.Call(opts, &out, "blockAndAggregate", calls)

	outstruct := new(struct {
		BlockNumber *big.Int
		BlockHash   [32]byte
		ReturnData  []Result
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BlockHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ReturnData = *abi.ConvertType(out[2], new([]Result)).(*[]Result)

	return *outstruct, err

}

func (_Multi *MultiCaller) TryAggregate(opts *bind.CallOpts, requireSuccess bool, calls []MultiCall) ([]Result, error) {
	var out []interface{}
	err := _Multi.contract.Call(opts, &out, "tryAggregate", requireSuccess, calls)

	if err != nil {
		return *new([]Result), err
	}

	out0 := *abi.ConvertType(out[0], new([]Result)).(*[]Result)

	return out0, err

}

func (_Multi *MultiCaller) TryBlockAndAggregate(opts *bind.CallOpts, requireSuccess bool, calls []MultiCall) (struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	ReturnData  []Result
}, error) {
	var out []interface{}
	err := _Multi.contract.Call(opts, &out, "tryBlockAndAggregate", requireSuccess, calls)

	outstruct := new(struct {
		BlockNumber *big.Int
		BlockHash   [32]byte
		ReturnData  []Result
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BlockHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ReturnData = *abi.ConvertType(out[2], new([]Result)).(*[]Result)

	return *outstruct, err

}