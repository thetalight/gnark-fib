package main

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"math/big"
	"time"
)

type Circuit struct {
	A, B   frontend.Variable
	Result frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	a := circuit.A
	b := circuit.B

	for i := 2; i <= 10000; i++ {
		next := api.Add(a, b)
		a, b = b, next
	}
	api.AssertIsEqual(circuit.Result, b)
	return nil
}

func main() {
	circuit, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("num of constraints: %d\n", circuit.GetNbConstraints())

	// Setup
	pk, vk, err := groth16.Setup(circuit)
	if err != nil {
		panic(err)
	}

	result, _ := new(big.Int).SetString("229563828224626036269062294847553432029308285020665226675966670650200296343", 10)
	assignment := &Circuit{A: 0, B: 1, Result: result}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}
	publicWitness, err := witness.Public()
	if err != nil {
		panic(err)
	}

	// Prove
	start := time.Now()
	proof, err := groth16.Prove(circuit, pk, witness)
	fmt.Printf("prover time: %s\n", time.Since(start))
	if err != nil {
		panic(err)
	}

	// Verify
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Verification Result: Success")
	}
}
