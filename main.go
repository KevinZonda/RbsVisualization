package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Operation struct {
	Name   string
	Inputs []string
	Output string
}

type Wire struct {
	Inputs string
	Output string
}

type Circuit struct {
	Operations []Operation
	Directions string
	Wiring     Wire
	Inputs     []string
}

func parseOp(line string) Operation {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return Operation{}
	}

	// Parse inputs from the domain field which is in format <.x,.y>
	domainStr := strings.Trim(fields[1], "<>")
	inputs := strings.Split(domainStr, ",")

	// Parse output which is in format .z
	output := strings.TrimSpace(fields[2])

	return Operation{
		Name:   fields[0],
		Inputs: inputs,
		Output: output,
	}
}

type IO struct {
}

func parseCircuit(input string) Circuit {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var circuit Circuit
	var currentSection string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines and separator lines
		if line == "" || strings.Contains(line, "---") {
			continue
		}

		// Parse header line
		if strings.Contains(line, "Name") {
			continue
		}

		// Check for section headers
		if strings.Contains(line, "Directions") {
			currentSection = "directions"
		} else if strings.Contains(line, "Wiring") {
			currentSection = "wiring"
		} else if strings.Contains(line, "Inputs") {
			currentSection = "inputs"
		}

		// Parse content based on current section
		ws := strings.Split(afterDash(line), "~")
		switch currentSection {
		case "":
			op := parseOp(line)
			if op.Name != "" {
				circuit.Operations = append(circuit.Operations, op)
			}
		case "directions":
			circuit.Directions = strings.Split(afterDash(line), "~")[0]
		case "wiring":
			circuit.Wiring = Wire{
				Inputs: ws[0],
			}
			if len(ws) == 2 {
				circuit.Wiring.Output = ws[1]
			}
		case "inputs":
			inputsStr := strings.TrimSpace(afterDash(line))
			inputsStr = strings.TrimLeft(inputsStr, "(")
			circuit.Inputs = strings.Fields(inputsStr)
		}
	}

	return circuit
}

func afterDash(s string) string {
	parts := strings.Split(s, "-")
	return parts[len(parts)-1]
}

func main() {
	rbs, _ := os.ReadFile("current.rbs")
	input := string(rbs)
	circuit := parseCircuit(input)

	// Print parsed circuit for verification
	fmt.Printf("Parsed %d operations\n", len(circuit.Operations))
	for _, op := range circuit.Operations {
		fmt.Printf("Operation: %s, Inputs: %v, Output: %s\n", op.Name, op.Inputs, op.Output)
	}
	fmt.Printf("Directions: %s\n", circuit.Directions)
	fmt.Printf("Wiring: %s\n", circuit.Wiring)
	fmt.Printf("Inputs: %v\n", circuit.Inputs)

	// Generate circuit visualization
	if err := visualizeCircuit(circuit); err != nil {
		fmt.Printf("Error generating visualization: %v\n", err)
		return
	}
	fmt.Println("Circuit visualization has been generated as 'circuit.png'")
}
