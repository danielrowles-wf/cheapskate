/*
Original work Copyright 2013 Eric Lesh
Modified work Copyright 2015 Tyler Treat

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
*/

package boom

import (
	"errors"
	"hash"
	"hash/fnv"
	"math"
)

var exp32 = math.Pow(2, 32)

// HyperLogLog implements the HyperLogLog cardinality estimation algorithm as
// described by Flajolet, Fusy, Gandouet, and Meunier in HyperLogLog: the
// analysis of a near-optimal cardinality estimation algorithm:
//
// http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf
//
// HyperLogLog is a probabilistic algorithm which approximates the number of
// distinct elements in a multiset. It works by hashing values and calculating
// the maximum number of leading zeros in the binary representation of each
// hash. If the maximum number of leading zeros is n, the estimated number of
// distinct elements in the set is 2^n. To minimize variance, the multiset is
// split into a configurable number of registers, the maximum number of leading
// zeros is calculated in the numbers in each register, and a harmonic mean is
// used to combine the estimates.
//
// For large or unbounded data sets, calculating the exact cardinality is
// impractical. HyperLogLog uses a fraction of the memory while providing an
// accurate approximation. For counting element frequency, refer to the
// Count-Min Sketch.
type HyperLogLog struct {
	registers []uint8     // counter registers
	m         uint        // number of registers
	b         uint32      // number of bits to calculate register
	alpha     float64     // bias-correction constant
	hash      hash.Hash32 // hash function
}

// NewHyperLogLog creates a new HyperLogLog with m registers. Returns an error
// if m isn't a power of two.
func NewHyperLogLog(m uint) (*HyperLogLog, error) {
	if (m & (m - 1)) != 0 {
		return nil, errors.New("m must be a power of two")
	}

	return &HyperLogLog{
		registers: make([]uint8, m),
		m:         m,
		b:         uint32(math.Ceil(math.Log2(float64(m)))),
		alpha:     calculateAlpha(m),
		hash:      fnv.New32(),
	}, nil
}

// NewDefaultHyperLogLog creates a new HyperLogLog optimized for the specified
// standard error. Returns an error if the number of registers can't be
// calculated for the provided accuracy.
func NewDefaultHyperLogLog(e float64) (*HyperLogLog, error) {
	m := math.Pow(1.04/e, 2)
	return NewHyperLogLog(uint(math.Pow(2, math.Ceil(math.Log2(m)))))
}

// Add will add the data to the set. Returns the HyperLogLog to allow for
// chaining.
func (h *HyperLogLog) Add(data []byte) *HyperLogLog {
	var (
		hash = h.calculateHash(data)
		k    = 32 - h.b
		r    = calculateRho(hash<<h.b, k)
		j    = hash >> uint(k)
	)

	if r > h.registers[j] {
		h.registers[j] = r
	}

	return h
}

// Count returns the approximated cardinality of the set.
func (h *HyperLogLog) Count() uint64 {
	sum := 0.0
	m := float64(h.m)
	for _, val := range h.registers {
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	estimate := h.alpha * m * m / sum
	if estimate <= 5.0/2.0*m {
		// Small range correction
		v := 0
		for _, r := range h.registers {
			if r == 0 {
				v++
			}
		}
		if v > 0 {
			estimate = m * math.Log(m/float64(v))
		}
	} else if estimate > 1.0/30.0*exp32 {
		// Large range correction
		estimate = -exp32 * math.Log(1-estimate/exp32)
	}
	return uint64(estimate)
}

// Merge combines this HyperLogLog with another. Returns an error if the number
// of registers in the two HyperLogLogs are not equal.
func (h *HyperLogLog) Merge(other *HyperLogLog) error {
	if h.m != other.m {
		return errors.New("number of registers must match")
	}

	for j, r := range other.registers {
		if r > h.registers[j] {
			h.registers[j] = r
		}
	}

	return nil
}

// Reset restores the HyperLogLog to its original state. It returns itself to
// allow for chaining.
func (h *HyperLogLog) Reset() *HyperLogLog {
	h.registers = make([]uint8, h.m)
	return h
}

// calculateHash calculates the 32-bit hash value for the provided data.
func (h *HyperLogLog) calculateHash(data []byte) uint32 {
	h.hash.Write(data)
	sum := h.hash.Sum32()
	h.hash.Reset()
	return sum
}

// calculateAlpha calculates the bias-correction constant alpha based on the
// number of registers, m.
func calculateAlpha(m uint) (result float64) {
	switch m {
	case 16:
		result = 0.673
	case 32:
		result = 0.697
	case 64:
		result = 0.709
	default:
		result = 0.7213 / (1.0 + 1.079/float64(m))
	}
	return result
}

// calculateRho calculates the position of the leftmost 1-bit.
func calculateRho(val, max uint32) uint8 {
	r := uint32(1)
	for val&0x80000000 == 0 && r <= max {
		r++
		val <<= 1
	}
	return uint8(r)
}
