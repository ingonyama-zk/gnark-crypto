// Copyright 2020 Consensys Software Inc.
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

package shplonk

import (
	"crypto/sha256"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/kzg"
)

// This example shows how to batch open a list of polynomials on a set of points,
// where each polynomial is opened on its own set of point.
// That is the i-th polynomial f_i is opened on  set of point S_i.
func Example_batchOpen() {

	const nbPolynomials = 10

	// sample a list of points and a list of polynomials. The i-th polynomial
	// is opened on the i-th set of points, there might be several points per set.
	points := make([][]fr.Element, nbPolynomials)
	polynomials := make([][]fr.Element, nbPolynomials)
	for i := 0; i < nbPolynomials; i++ {

		polynomials[i] = make([]fr.Element, 20+2*i) // random size
		for j := 0; j < 20+2*i; j++ {
			polynomials[i][j].SetRandom()
		}

		points[i] = make([]fr.Element, i+1) // random number of point
		for j := 0; j < i+1; j++ {
			points[i][j].SetRandom()
		}
	}

	// Create commitments for each polynomials
	var err error
	digests := make([]kzg.Digest, nbPolynomials)
	for i := 0; i < nbPolynomials; i++ {
		digests[i], err = kzg.Commit(polynomials[i], testSrs.Pk)
		if err != nil {
			panic(err)
		}
	}

	// hash function that is used for the challenge derivation in Fiat Shamir
	hf := sha256.New()

	// ceate an opening proof of polynomials[i] on the set points[i]
	openingProof, err := BatchOpen(polynomials, digests, points, hf, testSrs.Pk)
	if err != nil {
		panic(err)
	}

	// we verify the proof. If the proof is correct, then openingProof[i][j] contains
	// the evaluation of the polynomials[i] on points[i][j]
	err = BatchVerify(openingProof, digests, points, hf, testSrs.Vk)
	if err != nil {
		panic(err)
	}

	fmt.Println("verified")
	// output: verified
}
