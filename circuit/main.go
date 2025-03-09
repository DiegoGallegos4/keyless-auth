package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"

	"github.com/PolyhedraZK/ExpanderCompilerCollection/ecgo"
	"github.com/PolyhedraZK/ExpanderCompilerCollection/ecgo/test"
)

var (
	MAX_DEPTH = 256
)

// ZKAuthCircuit proves credential hash inclusion in a Merkle tree
type ZKAuthCircuit struct {
	Root          frontend.Variable   `gnark:",public"`
	ProofElements []frontend.Variable // private
	ProofIndex    frontend.Variable   // private
	Leaf          frontend.Variable   // private
}

// Define the zk circuit
func (circuit *ZKAuthCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	// Hash leaf
	h.Reset()
	h.Write(circuit.Leaf)
	hashed := h.Sum()

	depth := len(circuit.ProofElements)
	if depth == 0 {
		depth = MAX_DEPTH
	}
	proofIndices := api.ToBinary(circuit.ProofIndex, depth)

	// // Continuously hash with the proof elements
	for i := 0; i < len(circuit.ProofElements); i++ {
		element := circuit.ProofElements[i]
		// 0 = left, 1 = right
		index := proofIndices[i]

		d1 := api.Select(index, element, hashed)
		d2 := api.Select(index, hashed, element)

		h.Reset()
		h.Write(d1, d2)
		hashed = h.Sum()
	}

	// Verify calculates hash is equal to root
	api.AssertIsEqual(hashed, circuit.Root)
	return nil
}

func GenerateGroth16() error {
	var circuit ZKAuthCircuit

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return err
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return err
	}
	{
		rootDir, err := findRootDir()
		if err != nil {
			return fmt.Errorf("failed to find root directory: %w", err)
		}
		outputDir := filepath.Join(rootDir, "jsonK")
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		pubWitness := vk.NbPublicWitness()
		// Format and save verification key as JSON

		pw, err := json.Marshal(pubWitness)
		if err != nil {
			return fmt.Errorf("failed to marshal verification key to JSON: %w", err)
		}s

		err = os.WriteFile(filepath.Join(outputDir, "pw.json"), pw, 0644)
		if err != nil {
			return fmt.Errorf("failed to write verification key JSON: %w", err)
		}

		vkJSON, err := json.MarshalIndent(vk, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal verification key to JSON: %w", err)
		}

		err = os.WriteFile(filepath.Join(outputDir, "vk.json"), vkJSON, 0644)
		if err != nil {
			return fmt.Errorf("failed to write verification key JSON: %w", err)
		}

		proofJSON, err := json.MarshalIndent(pk, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal proof to JSON: %w", err)
		}

		err = os.WriteFile(filepath.Join(outputDir, "pk.json"), proofJSON, 0644)
		if err != nil {
			return fmt.Errorf("failed to write proof JSON: %w", err)
		}

	}
	{
		f, err := os.Create("mt.g16.vk")
		if err != nil {
			return err
		}
		_, err = vk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}
	{
		f, err := os.Create("mt.g16.pk")
		if err != nil {
			return err
		}
		_, err = pk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}
	{
		f, err := os.Create("contract_mt.g16.sol")
		if err != nil {
			return err
		}
		err = vk.ExportSolidity(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateExpander() error {
	assignment := &ZKAuthCircuit{
		Root: "123456789012345678901234567890123456789012345678901234567890abcd",
		ProofElements: []frontend.Variable{
			"234567890123456789012345678901234567890123456789012345678901dcba",
			"345678901234567890123456789012345678901234567890123456789012efab",
		},
		ProofIndex: 0,
		Leaf:       "123456789012345678901234567890123456789012345678901234567890fedc",
	}
	circuit, err := ecgo.Compile(ecc.BN254.ScalarField(), assignment)
	if err != nil {
		return err
	}

	c := circuit.GetLayeredCircuit()
	os.WriteFile("expander_circuit.txt", c.Serialize(), 0o644)
	inputSolver := circuit.GetInputSolver()
	witness, err := inputSolver.SolveInputAuto(nil)
	if err != nil {
		return err
	}
	os.WriteFile("expander_witness.txt", witness.Serialize(), 0o644)
	if !test.CheckCircuit(c, witness) {
		return errors.New("witness is not valid")
	}
	return nil
}

func main() {
	// Command line flags:
	// -expander: Generate expander circuit and witness files
	//   Example: ./main -expander
	//
	// -groth16: Generate Groth16 proving/verification keys and Solidity verifier
	//   Example: ./main -groth16
	expander := flag.Bool("expander", false, "Generate expander")
	groth16 := flag.Bool("groth16", false, "Generate Groth16 keys")
	flag.Parse()
	if !*expander && !*groth16 {
		log.Fatal("Please provide a command: 'expander' or 'groth16'")
	}

	if *expander {
		err := GenerateExpander()
		if err != nil {
			log.Fatalf("Failed to generate expander: %v", err)
		}
	} else if *groth16 {
		err := GenerateGroth16()
		if err != nil {
			log.Fatalf("Failed to generate Groth16 keys: %v", err)
		}
	} else {
		log.Fatalf("Invalid argument: %v", os.Args[1])
	}
}

func findRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check for common project root markers
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// We've reached the root of the file system
			return "", fmt.Errorf("could not find project root directory")
		}
		dir = parentDir
	}
}
