// Copyright 2016 The zerium Authors
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

package zrmclient

import "github.com/apolo-technologies/zerium"

// Verify that Client implements the abt interfaces.
var (
	_ = abt.ChainReader(&Client{})
	_ = abt.TransactionReader(&Client{})
	_ = abt.ChainStateReader(&Client{})
	_ = abt.ChainSyncReader(&Client{})
	_ = abt.ContractCaller(&Client{})
	_ = abt.GasEstimator(&Client{})
	_ = abt.GasPricer(&Client{})
	_ = abt.LogFilterer(&Client{})
	_ = abt.PendingStateReader(&Client{})
	// _ = abt.PendingStateEventer(&Client{})
	_ = abt.PendingContractCaller(&Client{})
)
