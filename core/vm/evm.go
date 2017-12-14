// Copyright 2017 The zerium Authors
// This file is part of the zerium library.
//
// The zerium library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The zerium library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the zerium library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"math/big"
	"sync/atomic"

	"github.com/apolo-technologies/zerium/zrmcom"
	"github.com/apolo-technologies/zerium/crypto"
	"github.com/apolo-technologies/zerium/params"
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

type (
	CanTransferFunc func(StateDB, zrmcom.Address, *big.Int) bool
	TransferFunc    func(StateDB, zrmcom.Address, zrmcom.Address, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH Zvm op code.
	GetHashFunc func(uint64) zrmcom.Hash
)

// run runs the given contract and takes care of running precompiles with a fallback to the byte code interpreter.
func run(zvm *Zvm, snapshot int, contract *Contract, input []byte) ([]byte, error) {
	if contract.CodeAddr != nil {
		precompiles := PrecompiledContractsHomestead
		if zvm.ChainConfig().IsByzantium(zvm.BlockNumber) {
			precompiles = PrecompiledContractsByzantium
		}
		if p := precompiles[*contract.CodeAddr]; p != nil {
			return RunPrecompiledContract(p, input, contract)
		}
	}
	return zvm.interpreter.Run(snapshot, contract, input)
}

// Context provides the Zvm with auxiliary information. Once provided
// it shouldn't be modified.
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Message information
	Origin   zrmcom.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase    zrmcom.Address // Provides information for COINBASE
	GasLimit    *big.Int       // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
}

// Zvm is the Zerium Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The Zvm should never be reused and is not thread safe.
type Zvm struct {
	// Context provides auxiliary blockchain related information
	Context
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// zvm.
	vmConfig Config
	// global (to this context) zerium virtual machine
	// used throughout the execution of the tx.
	interpreter *Interpreter
	// abort is used to abort the Zvm calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

// NewZvm retutrns a new Zvm . The returned Zvm is not thread safe and should
// only ever be used *once*.
func NewZvm(ctx Context, statedb StateDB, chainConfig *params.ChainConfig, vmConfig Config) *Zvm {
	zvm := &Zvm{
		Context:     ctx,
		StateDB:     statedb,
		vmConfig:    vmConfig,
		chainConfig: chainConfig,
		chainRules:  chainConfig.Rules(ctx.BlockNumber),
	}

	zvm.interpreter = NewInterpreter(zvm, vmConfig)
	return zvm
}

// Cancel cancels any running Zvm operation. This may be called concurrently and
// it's safe to be called multiple times.
func (zvm *Zvm) Cancel() {
	atomic.StoreInt32(&zvm.abort, 1)
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (zvm *Zvm) Call(caller ContractRef, addr zrmcom.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if zvm.vmConfig.NoRecursion && zvm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if zvm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if !zvm.Context.CanTransfer(zvm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}

	var (
		to       = AccountRef(addr)
		snapshot = zvm.StateDB.Snapshot()
	)
	if !zvm.StateDB.Exist(addr) {
		precompiles := PrecompiledContractsHomestead
		if zvm.ChainConfig().IsByzantium(zvm.BlockNumber) {
			precompiles = PrecompiledContractsByzantium
		}
		if precompiles[addr] == nil && zvm.ChainConfig().IsEIP158(zvm.BlockNumber) && value.Sign() == 0 {
			return nil, gas, nil
		}
		zvm.StateDB.CreateAccount(addr)
	}
	zvm.Transfer(zvm.StateDB, caller.Address(), to.Address(), value)

	// initialise a new contract and set the code that is to be used by the
	// E The contract is a scoped environment for this execution context
	// only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, zvm.StateDB.GetCodeHash(addr), zvm.StateDB.GetCode(addr))

	ret, err = run(zvm, snapshot, contract, input)
	// When an error was returned by the Zvm or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		zvm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (zvm *Zvm) CallCode(caller ContractRef, addr zrmcom.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if zvm.vmConfig.NoRecursion && zvm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if zvm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if !zvm.CanTransfer(zvm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}

	var (
		snapshot = zvm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)
	// initialise a new contract and set the code that is to be used by the
	// E The contract is a scoped zvmironment for this execution context
	// only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, zvm.StateDB.GetCodeHash(addr), zvm.StateDB.GetCode(addr))

	ret, err = run(zvm, snapshot, contract, input)
	if err != nil {
		zvm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (zvm *Zvm) DelegateCall(caller ContractRef, addr zrmcom.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if zvm.vmConfig.NoRecursion && zvm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if zvm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	var (
		snapshot = zvm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := NewContract(caller, to, nil, gas).AsDelegate()
	contract.SetCallCode(&addr, zvm.StateDB.GetCodeHash(addr), zvm.StateDB.GetCode(addr))

	ret, err = run(zvm, snapshot, contract, input)
	if err != nil {
		zvm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (zvm *Zvm) StaticCall(caller ContractRef, addr zrmcom.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if zvm.vmConfig.NoRecursion && zvm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if zvm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Make sure the readonly is only set if we aren't in readonly yet
	// this makes also sure that the readonly flag isn't removed for
	// child calls.
	if !zvm.interpreter.readOnly {
		zvm.interpreter.readOnly = true
		defer func() { zvm.interpreter.readOnly = false }()
	}

	var (
		to       = AccountRef(addr)
		snapshot = zvm.StateDB.Snapshot()
	)
	// Initialise a new contract and set the code that is to be used by the
	// Zvm. The contract is a scoped environment for this execution context
	// only.
	contract := NewContract(caller, to, new(big.Int), gas)
	contract.SetCallCode(&addr, zvm.StateDB.GetCodeHash(addr), zvm.StateDB.GetCode(addr))

	// When an error was returned by the Zvm or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in Homestead this also counts for code storage gas errors.
	ret, err = run(zvm, snapshot, contract, input)
	if err != nil {
		zvm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (zvm *Zvm) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr zrmcom.Address, leftOverGas uint64, err error) {

	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if zvm.depth > int(params.CallCreateDepth) {
		return nil, zrmcom.Address{}, gas, ErrDepth
	}
	if !zvm.CanTransfer(zvm.StateDB, caller.Address(), value) {
		return nil, zrmcom.Address{}, gas, ErrInsufficientBalance
	}
	// Ensure there's no existing contract already at the designated address
	nonce := zvm.StateDB.GetNonce(caller.Address())
	zvm.StateDB.SetNonce(caller.Address(), nonce+1)

	contractAddr = crypto.CreateAddress(caller.Address(), nonce)
	contractHash := zvm.StateDB.GetCodeHash(contractAddr)
	if zvm.StateDB.GetNonce(contractAddr) != 0 || (contractHash != (zrmcom.Hash{}) && contractHash != emptyCodeHash) {
		return nil, zrmcom.Address{}, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := zvm.StateDB.Snapshot()
	zvm.StateDB.CreateAccount(contractAddr)
	if zvm.ChainConfig().IsEIP158(zvm.BlockNumber) {
		zvm.StateDB.SetNonce(contractAddr, 1)
	}
	zvm.Transfer(zvm.StateDB, caller.Address(), contractAddr, value)

	// initialise a new contract and set the code that is to be used by the
	// E The contract is a scoped zvmironment for this execution context
	// only.
	contract := NewContract(caller, AccountRef(contractAddr), value, gas)
	contract.SetCallCode(&contractAddr, crypto.Keccak256Hash(code), code)

	if zvm.vmConfig.NoRecursion && zvm.depth > 0 {
		return nil, contractAddr, gas, nil
	}
	ret, err = run(zvm, snapshot, contract, nil)
	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := zvm.ChainConfig().IsEIP158(zvm.BlockNumber) && len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			zvm.StateDB.SetCode(contractAddr, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the Zvm or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (zvm.ChainConfig().IsHomestead(zvm.BlockNumber) || err != ErrCodeStoreOutOfGas)) {
		zvm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	return ret, contractAddr, contract.Gas, err
}

// ChainConfig returns the zvmironment's chain configuration
func (zvm *Zvm) ChainConfig() *params.ChainConfig { return zvm.chainConfig }

// Interpreter returns the Zvm interpreter
func (zvm *Zvm) Interpreter() *Interpreter { return zvm.interpreter }
