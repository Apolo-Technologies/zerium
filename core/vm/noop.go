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

func NoopCanTransfer(db StateDB, from zrmcom.Address, balance *big.Int) bool {
	return true
}
func NoopTransfer(db StateDB, from, to zrmcom.Address, amount *big.Int) {}

type NoopZvmCallContext struct{}

func (NoopZvmCallContext) Call(caller ContractRef, addr zrmcom.Address, data []byte, gas, value *big.Int) ([]byte, error) {
	return nil, nil
}
func (NoopZvmCallContext) CallCode(caller ContractRef, addr zrmcom.Address, data []byte, gas, value *big.Int) ([]byte, error) {
	return nil, nil
}
func (NoopZvmCallContext) Create(caller ContractRef, data []byte, gas, value *big.Int) ([]byte, zrmcom.Address, error) {
	return nil, zrmcom.Address{}, nil
}
func (NoopZvmCallContext) DelegateCall(me ContractRef, addr zrmcom.Address, data []byte, gas *big.Int) ([]byte, error) {
	return nil, nil
}

type NoopStateDB struct{}

func (NoopStateDB) CreateAccount(zrmcom.Address)                                       {}
func (NoopStateDB) SubBalance(zrmcom.Address, *big.Int)                                {}
func (NoopStateDB) AddBalance(zrmcom.Address, *big.Int)                                {}
func (NoopStateDB) GetBalance(zrmcom.Address) *big.Int                                 { return nil }
func (NoopStateDB) GetNonce(zrmcom.Address) uint64                                     { return 0 }
func (NoopStateDB) SetNonce(zrmcom.Address, uint64)                                    {}
func (NoopStateDB) GetCodeHash(zrmcom.Address) zrmcom.Hash                             { return zrmcom.Hash{} }
func (NoopStateDB) GetCode(zrmcom.Address) []byte                                      { return nil }
func (NoopStateDB) SetCode(zrmcom.Address, []byte)                                     {}
func (NoopStateDB) GetCodeSize(zrmcom.Address) int                                     { return 0 }
func (NoopStateDB) AddRefund(*big.Int)                                                 {}
func (NoopStateDB) GetRefund() *big.Int                                                { return nil }
func (NoopStateDB) GetState(zrmcom.Address, zrmcom.Hash) zrmcom.Hash                   { return zrmcom.Hash{} }
func (NoopStateDB) SetState(zrmcom.Address, zrmcom.Hash, zrmcom.Hash)                  {}
func (NoopStateDB) Suicide(zrmcom.Address) bool                                        { return false }
func (NoopStateDB) HasSuicided(zrmcom.Address) bool                                    { return false }
func (NoopStateDB) Exist(zrmcom.Address) bool                                          { return false }
func (NoopStateDB) Empty(zrmcom.Address) bool                                          { return false }
func (NoopStateDB) RevertToSnapshot(int)                                               {}
func (NoopStateDB) Snapshot() int                                                      { return 0 }
func (NoopStateDB) AddLog(*types.Log)                                                  {}
func (NoopStateDB) AddPreimage(zrmcom.Hash, []byte)                                    {}
func (NoopStateDB) ForEachStorage(zrmcom.Address, func(zrmcom.Hash, zrmcom.Hash) bool) {}
