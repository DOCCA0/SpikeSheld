// Code generated - DO NOT EDIT.
// This file is a simplified contract binding for InsurancePool

package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// InsurancePoolABI is the input ABI used to generate the binding from.
const InsurancePoolABI = `[{"inputs":[{"internalType":"address","name":"_usdt","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"OwnableInvalidOwner","type":"error"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"OwnableUnauthorizedAccount","type":"error"},{"inputs":[],"name":"ReentrancyGuardReentrantCall","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"newOracle","type":"address"}],"name":"OracleUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"user","type":"address"},{"indexed":false,"internalType":"uint256","name":"policyId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"PayoutExecuted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"user","type":"address"},{"indexed":false,"internalType":"uint256","name":"policyId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"premium","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"coverage","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"expiryTime","type":"uint256"}],"name":"PolicyPurchased","type":"event"},{"inputs":[],"name":"buyInsurance","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"coverageAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"coverageDuration","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"policyId","type":"uint256"}],"name":"executePayout","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"fundPool","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"policyId","type":"uint256"}],"name":"getPolicy","outputs":[{"components":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"premium","type":"uint256"},{"internalType":"uint256","name":"coverageAmount","type":"uint256"},{"internalType":"uint256","name":"purchaseTime","type":"uint256"},{"internalType":"uint256","name":"expiryTime","type":"uint256"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"bool","name":"claimed","type":"bool"}],"internalType":"struct InsurancePool.Policy","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getPoolBalance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getUserPoliciesCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"hasActivePolicy","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"oracle","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"premiumAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_oracle","type":"address"}],"name":"setOracle","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_premium","type":"uint256"},{"internalType":"uint256","name":"_coverage","type":"uint256"},{"internalType":"uint256","name":"_duration","type":"uint256"}],"name":"setPolicyParams","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"usdt","outputs":[{"internalType":"contract IERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"userPolicies","outputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"premium","type":"uint256"},{"internalType":"uint256","name":"coverageAmount","type":"uint256"},{"internalType":"uint256","name":"purchaseTime","type":"uint256"},{"internalType":"uint256","name":"expiryTime","type":"uint256"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"bool","name":"claimed","type":"bool"}],"stateMutability":"view","type":"function"}]`

// InsurancePool is an auto generated Go binding around an Ethereum contract.
type InsurancePool struct {
	InsurancePoolCaller     // Read-only binding to the contract
	InsurancePoolTransactor // Write-only binding to the contract
	InsurancePoolFilterer   // Log filterer for contract events
}

// InsurancePoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type InsurancePoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InsurancePoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InsurancePoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InsurancePoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InsurancePoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NewInsurancePool creates a new instance of InsurancePool, bound to a specific deployed contract.
func NewInsurancePool(address common.Address, backend bind.ContractBackend) (*InsurancePool, error) {
	contract, err := bindInsurancePool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InsurancePool{InsurancePoolCaller: InsurancePoolCaller{contract: contract}, InsurancePoolTransactor: InsurancePoolTransactor{contract: contract}, InsurancePoolFilterer: InsurancePoolFilterer{contract: contract}}, nil
}

// bindInsurancePool binds a generic wrapper to an already deployed contract.
func bindInsurancePool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(InsurancePoolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// ExecutePayout is a paid mutator transaction binding the contract method 0x9af1d35a.
//
// Solidity: function executePayout(address user, uint256 policyId) returns()
func (_InsurancePool *InsurancePoolTransactor) ExecutePayout(opts *bind.TransactOpts, user common.Address, policyId *big.Int) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "executePayout", user, policyId)
}

// ExecutePayout is a paid mutator transaction binding the contract method 0x9af1d35a.
//
// Solidity: function executePayout(address user, uint256 policyId) returns()
func (_InsurancePool *InsurancePool) ExecutePayout(opts *bind.TransactOpts, user common.Address, policyId *big.Int) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "executePayout", user, policyId)
}

// GetPoolBalance is a free data retrieval call binding the contract method 0x4550e706.
//
// Solidity: function getPoolBalance() view returns(uint256)
func (_InsurancePool *InsurancePoolCaller) GetPoolBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "getPoolBalance")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_InsurancePool *InsurancePoolCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "oracle")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}
