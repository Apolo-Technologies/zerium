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

package state

import (
	"encoding/json"
	"fmt"

	"github.com/apolo-technologies/zerium/zrmcom"
	"github.com/apolo-technologies/zerium/rlp"
	"github.com/apolo-technologies/zerium/trie"
)

type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

func (self *StateDB) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		obj := newObject(nil, zrmcom.BytesToAddress(addr), data, nil)
		account := DumpAccount{
			Balance:  data.Balance.String(),
			Nonce:    data.Nonce,
			Root:     zrmcom.Bytes2Hex(data.Root[:]),
			CodeHash: zrmcom.Bytes2Hex(data.CodeHash),
			Code:     zrmcom.Bytes2Hex(obj.Code(self.db)),
			Storage:  make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.Storage[zrmcom.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = zrmcom.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[zrmcom.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}
