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

// simple nonconcurrent reference implementation for hashsize segment based
// Binary Merkle tree hash on arbitrary but fixed maximum chunksize
//
// This implementation does not take advantage of any paralellisms and uses
// far more memory than necessary, but it is easy to see that it is correct.
// It can be used for generating test cases for optimized implementations.
// see testMerkleTreeHasherCorrectness function in merkletree_test.go
package merkletree

import (
	"hash"
)

// RefHasher is the non-optimized easy to read reference implementation of MerkleTree
type RefHasher struct {
	span    int
	section int
	cap     int
	h       hash.Hash
}

// NewRefHasher returns a new RefHasher
func NewRefHasher(hasher BaseHasher, count int) *RefHasher {
	h := hasher()
	hashsize := h.Size()
	maxsize := hashsize * count
	c := 2
	for ; c < count; c *= 2 {
	}
	if c > 2 {
		c /= 2
	}
	return &RefHasher{
		section: 2 * hashsize,
		span:    c * hashsize,
		cap:     maxsize,
		h:       h,
	}
}

// Hash returns the MerkleTree hash of the byte slice
// implements the SwarmHash interface
func (rh *RefHasher) Hash(d []byte) []byte {
	if len(d) > rh.cap {
		d = d[:rh.cap]
	}

	return rh.hash(d, rh.span)
}

func (rh *RefHasher) hash(d []byte, s int) []byte {
	l := len(d)
	left := d
	var right []byte
	if l > rh.section {
		for ; s >= l; s /= 2 {
		}
		left = rh.hash(d[:s], s)
		right = d[s:]
		if l-s > rh.section/2 {
			right = rh.hash(right, s)
		}
	}
	defer rh.h.Reset()
	rh.h.Write(left)
	rh.h.Write(right)
	h := rh.h.Sum(nil)
	return h
}