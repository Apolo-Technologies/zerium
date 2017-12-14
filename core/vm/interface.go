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

	"github.com/apolo-technologies/zerium/zrmcom"
	"github.com/apolo-technologies/zerium/core/types"
)

// StateDB is an Zvm database for full state querying.
type StateDB interface {
	CreateAccount(zrmcom.Address)

	SubBalance(zrmcom.Address, *big.Int)
	AddBalance(zrmcom.Address, *big.Int)
	GetBalance(zrmcom.Address) *big.Int

	GetNonce(zrmcom.Address) uint64
	SetNonce(zrmcom.Address, uint64)

	GetCodeHash(zrmcom.Address) zrmcom.Hash
	GetCode(zrmcom.Address) []byte
	SetCode(zrmcom.Address, []byte)
	GetCodeSize(zrmcom.Address) int

	AddRefund(*big.Int)
	GetRefund() *big.Int

	GetState(zrmcom.Address, zrmcom.Hash) zrmcom.Hash
	SetState(zrmcom.Address, zrmcom.Hash, zrmcom.Hash)

	Suicide(zrmcom.Address) bool
	HasSuicided(zrmcom.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(zrmcom.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(zrmcom.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(zrmcom.Hash, []byte)

	ForEachStorage(zrmcom.Address, func(zrmcom.Hash, zrmcom.Hash) bool)
}

// CallContext provides a basic interface for the Zvm calling conventions. The Zvm Zvm
// depends on this context being implemented for doing subcalls and initialising new Zvm contracts.
type CallContext interface {
	// Call another contract
	Call(env *Zvm, me ContractRef, addr zrmcom.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Take another's contract code and execute within our own context
	CallCode(env *Zvm, me ContractRef, addr zrmcom.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *Zvm, me ContractRef, addr zrmcom.Address, data []byte, gas *big.Int) ([]byte, error)
	// Create a new contract
	Create(env *Zvm, me ContractRef, data []byte, gas, value *big.Int) ([]byte, zrmcom.Address, error)
}
