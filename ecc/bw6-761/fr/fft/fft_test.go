// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fft

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bw6-761/fr"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestFFT(t *testing.T) {
	const maxSize = 1 << 10

	nbCosets := 3
	domainWithPrecompute := NewDomain(maxSize)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 5

	properties := gopter.NewProperties(parameters)

	properties.Property("DIF FFT should be consistent with dual basis", prop.ForAll(

		// checks that a random evaluation of a dual function eval(gen**ithpower) is consistent with the FFT result
		func(ithpower int) bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			domainWithPrecompute.FFT(pol, DIF)
			BitReverse(pol)

			sample := domainWithPrecompute.Generator
			sample.Exp(sample, big.NewInt(int64(ithpower)))

			eval := evaluatePolynomial(backupPol, sample)

			return eval.Equal(&pol[ithpower])

		},
		gen.IntRange(0, maxSize-1),
	))

	properties.Property("DIF FFT on cosets should be consistent with dual basis", prop.ForAll(

		// checks that a random evaluation of a dual function eval(gen**ithpower) is consistent with the FFT result
		func(ithpower int) bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			domainWithPrecompute.FFT(pol, DIF, WithCoset())
			BitReverse(pol)

			sample := domainWithPrecompute.Generator
			sample.Exp(sample, big.NewInt(int64(ithpower))).
				Mul(&sample, &domainWithPrecompute.FrMultiplicativeGen)

			eval := evaluatePolynomial(backupPol, sample)

			return eval.Equal(&pol[ithpower])

		},
		gen.IntRange(0, maxSize-1),
	))

	properties.Property("DIT FFT should be consistent with dual basis", prop.ForAll(

		// checks that a random evaluation of a dual function eval(gen**ithpower) is consistent with the FFT result
		func(ithpower int) bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			BitReverse(pol)
			domainWithPrecompute.FFT(pol, DIT)

			sample := domainWithPrecompute.Generator
			sample.Exp(sample, big.NewInt(int64(ithpower)))

			eval := evaluatePolynomial(backupPol, sample)

			return eval.Equal(&pol[ithpower])

		},
		gen.IntRange(0, maxSize-1),
	))

	properties.Property("bitReverse(DIF FFT(DIT FFT (bitReverse))))==id", prop.ForAll(

		func() bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			BitReverse(pol)
			domainWithPrecompute.FFT(pol, DIT)
			domainWithPrecompute.FFTInverse(pol, DIF)
			BitReverse(pol)

			check := true
			for i := 0; i < len(pol); i++ {
				check = check && pol[i].Equal(&backupPol[i])
			}
			return check
		},
	))

	properties.Property("bitReverse(DIF FFT(DIT FFT (bitReverse))))==id on cosets", prop.ForAll(

		func() bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			check := true

			for i := 1; i <= nbCosets; i++ {

				BitReverse(pol)
				domainWithPrecompute.FFT(pol, DIT, WithCoset())
				domainWithPrecompute.FFTInverse(pol, DIF, WithCoset())
				BitReverse(pol)

				for i := 0; i < len(pol); i++ {
					check = check && pol[i].Equal(&backupPol[i])
				}
			}

			return check
		},
	))

	properties.Property("DIT FFT(DIF FFT)==id", prop.ForAll(

		func() bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			domainWithPrecompute.FFTInverse(pol, DIF)
			domainWithPrecompute.FFT(pol, DIT)

			check := true
			for i := 0; i < len(pol); i++ {
				check = check && (pol[i] == backupPol[i])
			}
			return check
		},
	))

	properties.Property("DIT FFT(DIF FFT)==id on cosets", prop.ForAll(

		func() bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			domainWithPrecompute.FFTInverse(pol, DIF, WithCoset())
			domainWithPrecompute.FFT(pol, DIT, WithCoset())

			for i := 0; i < len(pol); i++ {
				if !(pol[i].Equal(&backupPol[i])) {
					return false
				}
			}

			// compute with nbTasks == 1
			domainWithPrecompute.FFTInverse(pol, DIF, WithCoset(), WithNbTasks(1))
			domainWithPrecompute.FFT(pol, DIT, WithCoset(), WithNbTasks(1))

			for i := 0; i < len(pol); i++ {
				if !(pol[i].Equal(&backupPol[i])) {
					return false
				}
			}

			return true
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// --------------------------------------------------------------------
// benches
func BenchmarkBitReverse(b *testing.B) {

	const maxSize = 1 << 20

	pol := make([]fr.Element, maxSize)
	pol[0].SetRandom()
	for i := 1; i < maxSize; i++ {
		pol[i] = pol[i-1]
	}

	for i := 8; i < 20; i++ {
		b.Run("bit reversing 2**"+strconv.Itoa(i)+"bits", func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				BitReverse(pol[:1<<i])
			}
		})
	}

}

func BenchmarkFFT(b *testing.B) {

	const maxSize = 1 << 20

	pol := make([]fr.Element, maxSize)
	pol[0].SetRandom()
	for i := 1; i < maxSize; i++ {
		pol[i] = pol[i-1]
	}

	for i := 8; i < 20; i++ {
		sizeDomain := 1 << i
		b.Run("fft 2**"+strconv.Itoa(i)+"bits", func(b *testing.B) {
			domain := NewDomain(uint64(sizeDomain))
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				domain.FFT(pol[:sizeDomain], DIT)
			}
		})
		b.Run("fft 2**"+strconv.Itoa(i)+"bits (coset)", func(b *testing.B) {
			domain := NewDomain(uint64(sizeDomain))
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				domain.FFT(pol[:sizeDomain], DIT, WithCoset())
			}
		})
	}

}

func BenchmarkFFTDITCosetReference(b *testing.B) {
	const maxSize = 1 << 20

	pol := make([]fr.Element, maxSize)
	pol[0].SetRandom()
	for i := 1; i < maxSize; i++ {
		pol[i] = pol[i-1]
	}

	domain := NewDomain(maxSize)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		domain.FFT(pol, DIT, WithCoset())
	}
}

func BenchmarkFFTDIFReference(b *testing.B) {
	const maxSize = 1 << 20

	pol := make([]fr.Element, maxSize)
	pol[0].SetRandom()
	for i := 1; i < maxSize; i++ {
		pol[i] = pol[i-1]
	}

	domain := NewDomain(maxSize)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		domain.FFT(pol, DIF)
	}
}

func evaluatePolynomial(pol []fr.Element, val fr.Element) fr.Element {
	var acc, res, tmp fr.Element
	res.Set(&pol[0])
	acc.Set(&val)
	for i := 1; i < len(pol); i++ {
		tmp.Mul(&acc, &pol[i])
		res.Add(&res, &tmp)
		acc.Mul(&acc, &val)
	}
	return res
}
