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

package runtime

import (
	"math/big"

	"github.com/apolo-technologies/zerium/zrmcom"
	"github.com/apolo-technologies/zerium/core"
	"github.com/apolo-technologies/zerium/core/vm"
)

func NewEnv(cfg *Config) *vm.Zvm {
	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) zrmcom.Hash { return zrmcom.Hash{} },

		Origin:      cfg.Origin,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    new(big.Int).SetUint64(cfg.GasLimit),
		GasPrice:    cfg.GasPrice,
	}

	return vm.NewZvm(context, cfg.State, cfg.ChainConfig, cfg.ZvmConfig)
}
