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

// Verify that Client implements the apolo-technologies interfaces.
var (
	_ = apolo-technologies.ChainReader(&Client{})
	_ = apolo-technologies.TransactionReader(&Client{})
	_ = apolo-technologies.ChainStateReader(&Client{})
	_ = apolo-technologies.ChainSyncReader(&Client{})
	_ = apolo-technologies.ContractCaller(&Client{})
	_ = apolo-technologies.GasEstimator(&Client{})
	_ = apolo-technologies.GasPricer(&Client{})
	_ = apolo-technologies.LogFilterer(&Client{})
	_ = apolo-technologies.PendingStateReader(&Client{})
	// _ = apolo-technologies.PendingStateEventer(&Client{})
	_ = apolo-technologies.PendingContractCaller(&Client{})
)
