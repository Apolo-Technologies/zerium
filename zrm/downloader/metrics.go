// Copyright 2015 The zerium Authors
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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/abt/zerium/metrics"
)

var (
	headerInMeter      = metrics.NewMeter("zrm/downloader/headers/in")
	headerReqTimer     = metrics.NewTimer("zrm/downloader/headers/req")
	headerDropMeter    = metrics.NewMeter("zrm/downloader/headers/drop")
	headerTimeoutMeter = metrics.NewMeter("zrm/downloader/headers/timeout")

	bodyInMeter      = metrics.NewMeter("zrm/downloader/bodies/in")
	bodyReqTimer     = metrics.NewTimer("zrm/downloader/bodies/req")
	bodyDropMeter    = metrics.NewMeter("zrm/downloader/bodies/drop")
	bodyTimeoutMeter = metrics.NewMeter("zrm/downloader/bodies/timeout")

	receiptInMeter      = metrics.NewMeter("zrm/downloader/receipts/in")
	receiptReqTimer     = metrics.NewTimer("zrm/downloader/receipts/req")
	receiptDropMeter    = metrics.NewMeter("zrm/downloader/receipts/drop")
	receiptTimeoutMeter = metrics.NewMeter("zrm/downloader/receipts/timeout")

	stateInMeter   = metrics.NewMeter("zrm/downloader/states/in")
	stateDropMeter = metrics.NewMeter("zrm/downloader/states/drop")
)
