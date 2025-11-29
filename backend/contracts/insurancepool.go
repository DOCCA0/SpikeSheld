// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
	_ = abi.ConvertType
)

// InsurancePoolPolicy is an auto generated low-level Go binding around an user-defined struct.
type InsurancePoolPolicy struct {
	User           common.Address
	Premium        *big.Int
	CoverageAmount *big.Int
	PurchaseTime   *big.Int
	ExpiryTime     *big.Int
	Active         bool
	Claimed        bool
}

// InsurancePoolMetaData contains all meta data concerning the InsurancePool contract.
var InsurancePoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOracle\",\"type\":\"address\"}],\"name\":\"OracleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"policyId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PayoutExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"policyId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"coverage\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiryTime\",\"type\":\"uint256\"}],\"name\":\"PolicyPurchased\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"buyInsurance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coverageAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coverageDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"policyId\",\"type\":\"uint256\"}],\"name\":\"executePayout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"policyId\",\"type\":\"uint256\"}],\"name\":\"getPolicy\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"coverageAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"purchaseTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"}],\"internalType\":\"structInsurancePool.Policy\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPoolBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserPolicies\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"coverageAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"purchaseTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"}],\"internalType\":\"structInsurancePool.Policy[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserPoliciesCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"hasActivePolicy\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_usdt\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"premiumAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_premium\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_coverage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_duration\",\"type\":\"uint256\"}],\"name\":\"setPolicyParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdt\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"userPolicies\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"coverageAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"purchaseTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// InsurancePoolABI is the input ABI used to generate the binding from.
// Deprecated: Use InsurancePoolMetaData.ABI instead.
var InsurancePoolABI = InsurancePoolMetaData.ABI

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

// InsurancePoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InsurancePoolSession struct {
	Contract     *InsurancePool    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InsurancePoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InsurancePoolCallerSession struct {
	Contract *InsurancePoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// InsurancePoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InsurancePoolTransactorSession struct {
	Contract     *InsurancePoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// InsurancePoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type InsurancePoolRaw struct {
	Contract *InsurancePool // Generic contract binding to access the raw methods on
}

// InsurancePoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InsurancePoolCallerRaw struct {
	Contract *InsurancePoolCaller // Generic read-only contract binding to access the raw methods on
}

// InsurancePoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InsurancePoolTransactorRaw struct {
	Contract *InsurancePoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInsurancePool creates a new instance of InsurancePool, bound to a specific deployed contract.
func NewInsurancePool(address common.Address, backend bind.ContractBackend) (*InsurancePool, error) {
	contract, err := bindInsurancePool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InsurancePool{InsurancePoolCaller: InsurancePoolCaller{contract: contract}, InsurancePoolTransactor: InsurancePoolTransactor{contract: contract}, InsurancePoolFilterer: InsurancePoolFilterer{contract: contract}}, nil
}

// NewInsurancePoolCaller creates a new read-only instance of InsurancePool, bound to a specific deployed contract.
func NewInsurancePoolCaller(address common.Address, caller bind.ContractCaller) (*InsurancePoolCaller, error) {
	contract, err := bindInsurancePool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolCaller{contract: contract}, nil
}

// NewInsurancePoolTransactor creates a new write-only instance of InsurancePool, bound to a specific deployed contract.
func NewInsurancePoolTransactor(address common.Address, transactor bind.ContractTransactor) (*InsurancePoolTransactor, error) {
	contract, err := bindInsurancePool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolTransactor{contract: contract}, nil
}

// NewInsurancePoolFilterer creates a new log filterer instance of InsurancePool, bound to a specific deployed contract.
func NewInsurancePoolFilterer(address common.Address, filterer bind.ContractFilterer) (*InsurancePoolFilterer, error) {
	contract, err := bindInsurancePool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolFilterer{contract: contract}, nil
}

// bindInsurancePool binds a generic wrapper to an already deployed contract.
func bindInsurancePool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InsurancePoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InsurancePool *InsurancePoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InsurancePool.Contract.InsurancePoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InsurancePool *InsurancePoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InsurancePool.Contract.InsurancePoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InsurancePool *InsurancePoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InsurancePool.Contract.InsurancePoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InsurancePool *InsurancePoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InsurancePool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InsurancePool *InsurancePoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InsurancePool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InsurancePool *InsurancePoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InsurancePool.Contract.contract.Transact(opts, method, params...)
}

// CoverageAmount is a free data retrieval call binding the contract method 0x67c0d00f.
//
// Solidity: function coverageAmount() view returns(uint256)
func (_InsurancePool *InsurancePoolCaller) CoverageAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "coverageAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CoverageAmount is a free data retrieval call binding the contract method 0x67c0d00f.
//
// Solidity: function coverageAmount() view returns(uint256)
func (_InsurancePool *InsurancePoolSession) CoverageAmount() (*big.Int, error) {
	return _InsurancePool.Contract.CoverageAmount(&_InsurancePool.CallOpts)
}

// CoverageAmount is a free data retrieval call binding the contract method 0x67c0d00f.
//
// Solidity: function coverageAmount() view returns(uint256)
func (_InsurancePool *InsurancePoolCallerSession) CoverageAmount() (*big.Int, error) {
	return _InsurancePool.Contract.CoverageAmount(&_InsurancePool.CallOpts)
}

// CoverageDuration is a free data retrieval call binding the contract method 0x39ff0616.
//
// Solidity: function coverageDuration() view returns(uint256)
func (_InsurancePool *InsurancePoolCaller) CoverageDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "coverageDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CoverageDuration is a free data retrieval call binding the contract method 0x39ff0616.
//
// Solidity: function coverageDuration() view returns(uint256)
func (_InsurancePool *InsurancePoolSession) CoverageDuration() (*big.Int, error) {
	return _InsurancePool.Contract.CoverageDuration(&_InsurancePool.CallOpts)
}

// CoverageDuration is a free data retrieval call binding the contract method 0x39ff0616.
//
// Solidity: function coverageDuration() view returns(uint256)
func (_InsurancePool *InsurancePoolCallerSession) CoverageDuration() (*big.Int, error) {
	return _InsurancePool.Contract.CoverageDuration(&_InsurancePool.CallOpts)
}

// GetPolicy is a free data retrieval call binding the contract method 0x66df322e.
//
// Solidity: function getPolicy(address user, uint256 policyId) view returns((address,uint256,uint256,uint256,uint256,bool,bool))
func (_InsurancePool *InsurancePoolCaller) GetPolicy(opts *bind.CallOpts, user common.Address, policyId *big.Int) (InsurancePoolPolicy, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "getPolicy", user, policyId)

	if err != nil {
		return *new(InsurancePoolPolicy), err
	}

	out0 := *abi.ConvertType(out[0], new(InsurancePoolPolicy)).(*InsurancePoolPolicy)

	return out0, err

}

// GetPolicy is a free data retrieval call binding the contract method 0x66df322e.
//
// Solidity: function getPolicy(address user, uint256 policyId) view returns((address,uint256,uint256,uint256,uint256,bool,bool))
func (_InsurancePool *InsurancePoolSession) GetPolicy(user common.Address, policyId *big.Int) (InsurancePoolPolicy, error) {
	return _InsurancePool.Contract.GetPolicy(&_InsurancePool.CallOpts, user, policyId)
}

// GetPolicy is a free data retrieval call binding the contract method 0x66df322e.
//
// Solidity: function getPolicy(address user, uint256 policyId) view returns((address,uint256,uint256,uint256,uint256,bool,bool))
func (_InsurancePool *InsurancePoolCallerSession) GetPolicy(user common.Address, policyId *big.Int) (InsurancePoolPolicy, error) {
	return _InsurancePool.Contract.GetPolicy(&_InsurancePool.CallOpts, user, policyId)
}

// GetPoolBalance is a free data retrieval call binding the contract method 0xabd70aa2.
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

// GetPoolBalance is a free data retrieval call binding the contract method 0xabd70aa2.
//
// Solidity: function getPoolBalance() view returns(uint256)
func (_InsurancePool *InsurancePoolSession) GetPoolBalance() (*big.Int, error) {
	return _InsurancePool.Contract.GetPoolBalance(&_InsurancePool.CallOpts)
}

// GetPoolBalance is a free data retrieval call binding the contract method 0xabd70aa2.
//
// Solidity: function getPoolBalance() view returns(uint256)
func (_InsurancePool *InsurancePoolCallerSession) GetPoolBalance() (*big.Int, error) {
	return _InsurancePool.Contract.GetPoolBalance(&_InsurancePool.CallOpts)
}

// GetUserPolicies is a free data retrieval call binding the contract method 0x19ac4614.
//
// Solidity: function getUserPolicies(address user) view returns((address,uint256,uint256,uint256,uint256,bool,bool)[])
func (_InsurancePool *InsurancePoolCaller) GetUserPolicies(opts *bind.CallOpts, user common.Address) ([]InsurancePoolPolicy, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "getUserPolicies", user)

	if err != nil {
		return *new([]InsurancePoolPolicy), err
	}

	out0 := *abi.ConvertType(out[0], new([]InsurancePoolPolicy)).(*[]InsurancePoolPolicy)

	return out0, err

}

// GetUserPolicies is a free data retrieval call binding the contract method 0x19ac4614.
//
// Solidity: function getUserPolicies(address user) view returns((address,uint256,uint256,uint256,uint256,bool,bool)[])
func (_InsurancePool *InsurancePoolSession) GetUserPolicies(user common.Address) ([]InsurancePoolPolicy, error) {
	return _InsurancePool.Contract.GetUserPolicies(&_InsurancePool.CallOpts, user)
}

// GetUserPolicies is a free data retrieval call binding the contract method 0x19ac4614.
//
// Solidity: function getUserPolicies(address user) view returns((address,uint256,uint256,uint256,uint256,bool,bool)[])
func (_InsurancePool *InsurancePoolCallerSession) GetUserPolicies(user common.Address) ([]InsurancePoolPolicy, error) {
	return _InsurancePool.Contract.GetUserPolicies(&_InsurancePool.CallOpts, user)
}

// GetUserPoliciesCount is a free data retrieval call binding the contract method 0x610134e4.
//
// Solidity: function getUserPoliciesCount(address user) view returns(uint256)
func (_InsurancePool *InsurancePoolCaller) GetUserPoliciesCount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "getUserPoliciesCount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserPoliciesCount is a free data retrieval call binding the contract method 0x610134e4.
//
// Solidity: function getUserPoliciesCount(address user) view returns(uint256)
func (_InsurancePool *InsurancePoolSession) GetUserPoliciesCount(user common.Address) (*big.Int, error) {
	return _InsurancePool.Contract.GetUserPoliciesCount(&_InsurancePool.CallOpts, user)
}

// GetUserPoliciesCount is a free data retrieval call binding the contract method 0x610134e4.
//
// Solidity: function getUserPoliciesCount(address user) view returns(uint256)
func (_InsurancePool *InsurancePoolCallerSession) GetUserPoliciesCount(user common.Address) (*big.Int, error) {
	return _InsurancePool.Contract.GetUserPoliciesCount(&_InsurancePool.CallOpts, user)
}

// HasActivePolicy is a free data retrieval call binding the contract method 0xb0d5e4d5.
//
// Solidity: function hasActivePolicy(address user) view returns(bool)
func (_InsurancePool *InsurancePoolCaller) HasActivePolicy(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "hasActivePolicy", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasActivePolicy is a free data retrieval call binding the contract method 0xb0d5e4d5.
//
// Solidity: function hasActivePolicy(address user) view returns(bool)
func (_InsurancePool *InsurancePoolSession) HasActivePolicy(user common.Address) (bool, error) {
	return _InsurancePool.Contract.HasActivePolicy(&_InsurancePool.CallOpts, user)
}

// HasActivePolicy is a free data retrieval call binding the contract method 0xb0d5e4d5.
//
// Solidity: function hasActivePolicy(address user) view returns(bool)
func (_InsurancePool *InsurancePoolCallerSession) HasActivePolicy(user common.Address) (bool, error) {
	return _InsurancePool.Contract.HasActivePolicy(&_InsurancePool.CallOpts, user)
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

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_InsurancePool *InsurancePoolSession) Oracle() (common.Address, error) {
	return _InsurancePool.Contract.Oracle(&_InsurancePool.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_InsurancePool *InsurancePoolCallerSession) Oracle() (common.Address, error) {
	return _InsurancePool.Contract.Oracle(&_InsurancePool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_InsurancePool *InsurancePoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_InsurancePool *InsurancePoolSession) Owner() (common.Address, error) {
	return _InsurancePool.Contract.Owner(&_InsurancePool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_InsurancePool *InsurancePoolCallerSession) Owner() (common.Address, error) {
	return _InsurancePool.Contract.Owner(&_InsurancePool.CallOpts)
}

// PremiumAmount is a free data retrieval call binding the contract method 0x44530f3a.
//
// Solidity: function premiumAmount() view returns(uint256)
func (_InsurancePool *InsurancePoolCaller) PremiumAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "premiumAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PremiumAmount is a free data retrieval call binding the contract method 0x44530f3a.
//
// Solidity: function premiumAmount() view returns(uint256)
func (_InsurancePool *InsurancePoolSession) PremiumAmount() (*big.Int, error) {
	return _InsurancePool.Contract.PremiumAmount(&_InsurancePool.CallOpts)
}

// PremiumAmount is a free data retrieval call binding the contract method 0x44530f3a.
//
// Solidity: function premiumAmount() view returns(uint256)
func (_InsurancePool *InsurancePoolCallerSession) PremiumAmount() (*big.Int, error) {
	return _InsurancePool.Contract.PremiumAmount(&_InsurancePool.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_InsurancePool *InsurancePoolCaller) Usdt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "usdt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_InsurancePool *InsurancePoolSession) Usdt() (common.Address, error) {
	return _InsurancePool.Contract.Usdt(&_InsurancePool.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_InsurancePool *InsurancePoolCallerSession) Usdt() (common.Address, error) {
	return _InsurancePool.Contract.Usdt(&_InsurancePool.CallOpts)
}

// UserPolicies is a free data retrieval call binding the contract method 0x3d36adc5.
//
// Solidity: function userPolicies(address , uint256 ) view returns(address user, uint256 premium, uint256 coverageAmount, uint256 purchaseTime, uint256 expiryTime, bool active, bool claimed)
func (_InsurancePool *InsurancePoolCaller) UserPolicies(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	User           common.Address
	Premium        *big.Int
	CoverageAmount *big.Int
	PurchaseTime   *big.Int
	ExpiryTime     *big.Int
	Active         bool
	Claimed        bool
}, error) {
	var out []interface{}
	err := _InsurancePool.contract.Call(opts, &out, "userPolicies", arg0, arg1)

	outstruct := new(struct {
		User           common.Address
		Premium        *big.Int
		CoverageAmount *big.Int
		PurchaseTime   *big.Int
		ExpiryTime     *big.Int
		Active         bool
		Claimed        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Premium = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.CoverageAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PurchaseTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.ExpiryTime = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.Claimed = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// UserPolicies is a free data retrieval call binding the contract method 0x3d36adc5.
//
// Solidity: function userPolicies(address , uint256 ) view returns(address user, uint256 premium, uint256 coverageAmount, uint256 purchaseTime, uint256 expiryTime, bool active, bool claimed)
func (_InsurancePool *InsurancePoolSession) UserPolicies(arg0 common.Address, arg1 *big.Int) (struct {
	User           common.Address
	Premium        *big.Int
	CoverageAmount *big.Int
	PurchaseTime   *big.Int
	ExpiryTime     *big.Int
	Active         bool
	Claimed        bool
}, error) {
	return _InsurancePool.Contract.UserPolicies(&_InsurancePool.CallOpts, arg0, arg1)
}

// UserPolicies is a free data retrieval call binding the contract method 0x3d36adc5.
//
// Solidity: function userPolicies(address , uint256 ) view returns(address user, uint256 premium, uint256 coverageAmount, uint256 purchaseTime, uint256 expiryTime, bool active, bool claimed)
func (_InsurancePool *InsurancePoolCallerSession) UserPolicies(arg0 common.Address, arg1 *big.Int) (struct {
	User           common.Address
	Premium        *big.Int
	CoverageAmount *big.Int
	PurchaseTime   *big.Int
	ExpiryTime     *big.Int
	Active         bool
	Claimed        bool
}, error) {
	return _InsurancePool.Contract.UserPolicies(&_InsurancePool.CallOpts, arg0, arg1)
}

// BuyInsurance is a paid mutator transaction binding the contract method 0x71618c79.
//
// Solidity: function buyInsurance() returns()
func (_InsurancePool *InsurancePoolTransactor) BuyInsurance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "buyInsurance")
}

// BuyInsurance is a paid mutator transaction binding the contract method 0x71618c79.
//
// Solidity: function buyInsurance() returns()
func (_InsurancePool *InsurancePoolSession) BuyInsurance() (*types.Transaction, error) {
	return _InsurancePool.Contract.BuyInsurance(&_InsurancePool.TransactOpts)
}

// BuyInsurance is a paid mutator transaction binding the contract method 0x71618c79.
//
// Solidity: function buyInsurance() returns()
func (_InsurancePool *InsurancePoolTransactorSession) BuyInsurance() (*types.Transaction, error) {
	return _InsurancePool.Contract.BuyInsurance(&_InsurancePool.TransactOpts)
}

// ExecutePayout is a paid mutator transaction binding the contract method 0x1ccd0b2e.
//
// Solidity: function executePayout(address user, uint256 policyId) returns()
func (_InsurancePool *InsurancePoolTransactor) ExecutePayout(opts *bind.TransactOpts, user common.Address, policyId *big.Int) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "executePayout", user, policyId)
}

// ExecutePayout is a paid mutator transaction binding the contract method 0x1ccd0b2e.
//
// Solidity: function executePayout(address user, uint256 policyId) returns()
func (_InsurancePool *InsurancePoolSession) ExecutePayout(user common.Address, policyId *big.Int) (*types.Transaction, error) {
	return _InsurancePool.Contract.ExecutePayout(&_InsurancePool.TransactOpts, user, policyId)
}

// ExecutePayout is a paid mutator transaction binding the contract method 0x1ccd0b2e.
//
// Solidity: function executePayout(address user, uint256 policyId) returns()
func (_InsurancePool *InsurancePoolTransactorSession) ExecutePayout(user common.Address, policyId *big.Int) (*types.Transaction, error) {
	return _InsurancePool.Contract.ExecutePayout(&_InsurancePool.TransactOpts, user, policyId)
}

// FundPool is a paid mutator transaction binding the contract method 0xda3a5f72.
//
// Solidity: function fundPool(uint256 amount) returns()
func (_InsurancePool *InsurancePoolTransactor) FundPool(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "fundPool", amount)
}

// FundPool is a paid mutator transaction binding the contract method 0xda3a5f72.
//
// Solidity: function fundPool(uint256 amount) returns()
func (_InsurancePool *InsurancePoolSession) FundPool(amount *big.Int) (*types.Transaction, error) {
	return _InsurancePool.Contract.FundPool(&_InsurancePool.TransactOpts, amount)
}

// FundPool is a paid mutator transaction binding the contract method 0xda3a5f72.
//
// Solidity: function fundPool(uint256 amount) returns()
func (_InsurancePool *InsurancePoolTransactorSession) FundPool(amount *big.Int) (*types.Transaction, error) {
	return _InsurancePool.Contract.FundPool(&_InsurancePool.TransactOpts, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _usdt) returns()
func (_InsurancePool *InsurancePoolTransactor) Initialize(opts *bind.TransactOpts, _usdt common.Address) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "initialize", _usdt)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _usdt) returns()
func (_InsurancePool *InsurancePoolSession) Initialize(_usdt common.Address) (*types.Transaction, error) {
	return _InsurancePool.Contract.Initialize(&_InsurancePool.TransactOpts, _usdt)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _usdt) returns()
func (_InsurancePool *InsurancePoolTransactorSession) Initialize(_usdt common.Address) (*types.Transaction, error) {
	return _InsurancePool.Contract.Initialize(&_InsurancePool.TransactOpts, _usdt)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_InsurancePool *InsurancePoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_InsurancePool *InsurancePoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _InsurancePool.Contract.RenounceOwnership(&_InsurancePool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_InsurancePool *InsurancePoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _InsurancePool.Contract.RenounceOwnership(&_InsurancePool.TransactOpts)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_InsurancePool *InsurancePoolTransactor) SetOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "setOracle", _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_InsurancePool *InsurancePoolSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _InsurancePool.Contract.SetOracle(&_InsurancePool.TransactOpts, _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_InsurancePool *InsurancePoolTransactorSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _InsurancePool.Contract.SetOracle(&_InsurancePool.TransactOpts, _oracle)
}

// SetPolicyParams is a paid mutator transaction binding the contract method 0xd14e6407.
//
// Solidity: function setPolicyParams(uint256 _premium, uint256 _coverage, uint256 _duration) returns()
func (_InsurancePool *InsurancePoolTransactor) SetPolicyParams(opts *bind.TransactOpts, _premium *big.Int, _coverage *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "setPolicyParams", _premium, _coverage, _duration)
}

// SetPolicyParams is a paid mutator transaction binding the contract method 0xd14e6407.
//
// Solidity: function setPolicyParams(uint256 _premium, uint256 _coverage, uint256 _duration) returns()
func (_InsurancePool *InsurancePoolSession) SetPolicyParams(_premium *big.Int, _coverage *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _InsurancePool.Contract.SetPolicyParams(&_InsurancePool.TransactOpts, _premium, _coverage, _duration)
}

// SetPolicyParams is a paid mutator transaction binding the contract method 0xd14e6407.
//
// Solidity: function setPolicyParams(uint256 _premium, uint256 _coverage, uint256 _duration) returns()
func (_InsurancePool *InsurancePoolTransactorSession) SetPolicyParams(_premium *big.Int, _coverage *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _InsurancePool.Contract.SetPolicyParams(&_InsurancePool.TransactOpts, _premium, _coverage, _duration)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_InsurancePool *InsurancePoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _InsurancePool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_InsurancePool *InsurancePoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _InsurancePool.Contract.TransferOwnership(&_InsurancePool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_InsurancePool *InsurancePoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _InsurancePool.Contract.TransferOwnership(&_InsurancePool.TransactOpts, newOwner)
}

// InsurancePoolInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the InsurancePool contract.
type InsurancePoolInitializedIterator struct {
	Event *InsurancePoolInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *InsurancePoolInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InsurancePoolInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(InsurancePoolInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *InsurancePoolInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InsurancePoolInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InsurancePoolInitialized represents a Initialized event raised by the InsurancePool contract.
type InsurancePoolInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_InsurancePool *InsurancePoolFilterer) FilterInitialized(opts *bind.FilterOpts) (*InsurancePoolInitializedIterator, error) {

	logs, sub, err := _InsurancePool.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &InsurancePoolInitializedIterator{contract: _InsurancePool.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_InsurancePool *InsurancePoolFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *InsurancePoolInitialized) (event.Subscription, error) {

	logs, sub, err := _InsurancePool.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InsurancePoolInitialized)
				if err := _InsurancePool.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_InsurancePool *InsurancePoolFilterer) ParseInitialized(log types.Log) (*InsurancePoolInitialized, error) {
	event := new(InsurancePoolInitialized)
	if err := _InsurancePool.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InsurancePoolOracleUpdatedIterator is returned from FilterOracleUpdated and is used to iterate over the raw logs and unpacked data for OracleUpdated events raised by the InsurancePool contract.
type InsurancePoolOracleUpdatedIterator struct {
	Event *InsurancePoolOracleUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *InsurancePoolOracleUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InsurancePoolOracleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(InsurancePoolOracleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *InsurancePoolOracleUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InsurancePoolOracleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InsurancePoolOracleUpdated represents a OracleUpdated event raised by the InsurancePool contract.
type InsurancePoolOracleUpdated struct {
	NewOracle common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleUpdated is a free log retrieval operation binding the contract event 0x3df77beb5db05fcdd70a30fc8adf3f83f9501b68579455adbd100b8180940394.
//
// Solidity: event OracleUpdated(address indexed newOracle)
func (_InsurancePool *InsurancePoolFilterer) FilterOracleUpdated(opts *bind.FilterOpts, newOracle []common.Address) (*InsurancePoolOracleUpdatedIterator, error) {

	var newOracleRule []interface{}
	for _, newOracleItem := range newOracle {
		newOracleRule = append(newOracleRule, newOracleItem)
	}

	logs, sub, err := _InsurancePool.contract.FilterLogs(opts, "OracleUpdated", newOracleRule)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolOracleUpdatedIterator{contract: _InsurancePool.contract, event: "OracleUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleUpdated is a free log subscription operation binding the contract event 0x3df77beb5db05fcdd70a30fc8adf3f83f9501b68579455adbd100b8180940394.
//
// Solidity: event OracleUpdated(address indexed newOracle)
func (_InsurancePool *InsurancePoolFilterer) WatchOracleUpdated(opts *bind.WatchOpts, sink chan<- *InsurancePoolOracleUpdated, newOracle []common.Address) (event.Subscription, error) {

	var newOracleRule []interface{}
	for _, newOracleItem := range newOracle {
		newOracleRule = append(newOracleRule, newOracleItem)
	}

	logs, sub, err := _InsurancePool.contract.WatchLogs(opts, "OracleUpdated", newOracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InsurancePoolOracleUpdated)
				if err := _InsurancePool.contract.UnpackLog(event, "OracleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOracleUpdated is a log parse operation binding the contract event 0x3df77beb5db05fcdd70a30fc8adf3f83f9501b68579455adbd100b8180940394.
//
// Solidity: event OracleUpdated(address indexed newOracle)
func (_InsurancePool *InsurancePoolFilterer) ParseOracleUpdated(log types.Log) (*InsurancePoolOracleUpdated, error) {
	event := new(InsurancePoolOracleUpdated)
	if err := _InsurancePool.contract.UnpackLog(event, "OracleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InsurancePoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the InsurancePool contract.
type InsurancePoolOwnershipTransferredIterator struct {
	Event *InsurancePoolOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *InsurancePoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InsurancePoolOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(InsurancePoolOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *InsurancePoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InsurancePoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InsurancePoolOwnershipTransferred represents a OwnershipTransferred event raised by the InsurancePool contract.
type InsurancePoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_InsurancePool *InsurancePoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*InsurancePoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _InsurancePool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolOwnershipTransferredIterator{contract: _InsurancePool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_InsurancePool *InsurancePoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *InsurancePoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _InsurancePool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InsurancePoolOwnershipTransferred)
				if err := _InsurancePool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_InsurancePool *InsurancePoolFilterer) ParseOwnershipTransferred(log types.Log) (*InsurancePoolOwnershipTransferred, error) {
	event := new(InsurancePoolOwnershipTransferred)
	if err := _InsurancePool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InsurancePoolPayoutExecutedIterator is returned from FilterPayoutExecuted and is used to iterate over the raw logs and unpacked data for PayoutExecuted events raised by the InsurancePool contract.
type InsurancePoolPayoutExecutedIterator struct {
	Event *InsurancePoolPayoutExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *InsurancePoolPayoutExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InsurancePoolPayoutExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(InsurancePoolPayoutExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *InsurancePoolPayoutExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InsurancePoolPayoutExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InsurancePoolPayoutExecuted represents a PayoutExecuted event raised by the InsurancePool contract.
type InsurancePoolPayoutExecuted struct {
	User     common.Address
	PolicyId *big.Int
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPayoutExecuted is a free log retrieval operation binding the contract event 0x63f55f28cc9b6363ecdbaa936809f86425d5709291a6d3e58910f0d9fc4b390f.
//
// Solidity: event PayoutExecuted(address indexed user, uint256 policyId, uint256 amount)
func (_InsurancePool *InsurancePoolFilterer) FilterPayoutExecuted(opts *bind.FilterOpts, user []common.Address) (*InsurancePoolPayoutExecutedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _InsurancePool.contract.FilterLogs(opts, "PayoutExecuted", userRule)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolPayoutExecutedIterator{contract: _InsurancePool.contract, event: "PayoutExecuted", logs: logs, sub: sub}, nil
}

// WatchPayoutExecuted is a free log subscription operation binding the contract event 0x63f55f28cc9b6363ecdbaa936809f86425d5709291a6d3e58910f0d9fc4b390f.
//
// Solidity: event PayoutExecuted(address indexed user, uint256 policyId, uint256 amount)
func (_InsurancePool *InsurancePoolFilterer) WatchPayoutExecuted(opts *bind.WatchOpts, sink chan<- *InsurancePoolPayoutExecuted, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _InsurancePool.contract.WatchLogs(opts, "PayoutExecuted", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InsurancePoolPayoutExecuted)
				if err := _InsurancePool.contract.UnpackLog(event, "PayoutExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayoutExecuted is a log parse operation binding the contract event 0x63f55f28cc9b6363ecdbaa936809f86425d5709291a6d3e58910f0d9fc4b390f.
//
// Solidity: event PayoutExecuted(address indexed user, uint256 policyId, uint256 amount)
func (_InsurancePool *InsurancePoolFilterer) ParsePayoutExecuted(log types.Log) (*InsurancePoolPayoutExecuted, error) {
	event := new(InsurancePoolPayoutExecuted)
	if err := _InsurancePool.contract.UnpackLog(event, "PayoutExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InsurancePoolPolicyPurchasedIterator is returned from FilterPolicyPurchased and is used to iterate over the raw logs and unpacked data for PolicyPurchased events raised by the InsurancePool contract.
type InsurancePoolPolicyPurchasedIterator struct {
	Event *InsurancePoolPolicyPurchased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *InsurancePoolPolicyPurchasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InsurancePoolPolicyPurchased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(InsurancePoolPolicyPurchased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *InsurancePoolPolicyPurchasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InsurancePoolPolicyPurchasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InsurancePoolPolicyPurchased represents a PolicyPurchased event raised by the InsurancePool contract.
type InsurancePoolPolicyPurchased struct {
	User       common.Address
	PolicyId   *big.Int
	Premium    *big.Int
	Coverage   *big.Int
	ExpiryTime *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPolicyPurchased is a free log retrieval operation binding the contract event 0x8e63e251f17c4f8c8dd082b9641ba8a3a04d5506e4a1b1782c6e3f57abc39374.
//
// Solidity: event PolicyPurchased(address indexed user, uint256 policyId, uint256 premium, uint256 coverage, uint256 expiryTime)
func (_InsurancePool *InsurancePoolFilterer) FilterPolicyPurchased(opts *bind.FilterOpts, user []common.Address) (*InsurancePoolPolicyPurchasedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _InsurancePool.contract.FilterLogs(opts, "PolicyPurchased", userRule)
	if err != nil {
		return nil, err
	}
	return &InsurancePoolPolicyPurchasedIterator{contract: _InsurancePool.contract, event: "PolicyPurchased", logs: logs, sub: sub}, nil
}

// WatchPolicyPurchased is a free log subscription operation binding the contract event 0x8e63e251f17c4f8c8dd082b9641ba8a3a04d5506e4a1b1782c6e3f57abc39374.
//
// Solidity: event PolicyPurchased(address indexed user, uint256 policyId, uint256 premium, uint256 coverage, uint256 expiryTime)
func (_InsurancePool *InsurancePoolFilterer) WatchPolicyPurchased(opts *bind.WatchOpts, sink chan<- *InsurancePoolPolicyPurchased, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _InsurancePool.contract.WatchLogs(opts, "PolicyPurchased", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InsurancePoolPolicyPurchased)
				if err := _InsurancePool.contract.UnpackLog(event, "PolicyPurchased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePolicyPurchased is a log parse operation binding the contract event 0x8e63e251f17c4f8c8dd082b9641ba8a3a04d5506e4a1b1782c6e3f57abc39374.
//
// Solidity: event PolicyPurchased(address indexed user, uint256 policyId, uint256 premium, uint256 coverage, uint256 expiryTime)
func (_InsurancePool *InsurancePoolFilterer) ParsePolicyPurchased(log types.Log) (*InsurancePoolPolicyPurchased, error) {
	event := new(InsurancePoolPolicyPurchased)
	if err := _InsurancePool.contract.UnpackLog(event, "PolicyPurchased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
