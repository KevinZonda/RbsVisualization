package main

import (
	"fmt"
	"os"
	"os/exec"
)

func generateDotFile(circuit Circuit) string {
	dot := "digraph circuit {\n"
	dot += "\tnode [shape=box];\n"

	// Add operation nodes
	for i, op := range circuit.Operations {
		dot += fmt.Sprintf("\top%d [label=\"%s\"];\n", i, op.Name)
		// Add edges from inputs to operation
		for _, input := range op.Inputs {
			dot += fmt.Sprintf("\t\"%s\" -> op%d;\n", input, i)
		}
		// Add edge from operation to output
		dot += fmt.Sprintf("\top%d -> \"%s\";\n", i, op.Output)
	}

	// Add input nodes with different shape
	dot += "\tnode [shape=ellipse];\n"
	for _, input := range circuit.Inputs {
		dot += fmt.Sprintf("\t\"%s\" [style=filled,fillcolor=lightblue];\n", input)
	}

	dot += "}\n"
	return dot
}

func visualizeCircuit(circuit Circuit) error {
	dotContent := generateDotFile(circuit)

	if err := os.WriteFile("circuit.dot", []byte(dotContent), 0644); err != nil {
		return fmt.Errorf("failed to write DOT file: %v", err)
	}

	cmd := exec.Command("dot", "-Tpng", "-o", "circuit.png", "circuit.dot")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate PNG: %v", err)
	}

	return nil
}
