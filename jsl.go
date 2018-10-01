package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {

	fmt.Println("JSL")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	programStack := &stack{make([]langObject, 0)}
	programSymbolTable := &symbolTable{make(map[uint64]*symbolTableEntry),}
	programVariableScope := &variableScope{make(map[string]*langVariable),nil,}

	for true {

		fmt.Print("> ")

		if !scanner.Scan() {
			fmt.Print("\n")
			os.Exit(0)
		}

		input := scanner.Text()

		if input == "exit" || input == "quit" {
			os.Exit(0)
		}

		err := evalString(input, programStack, programVariableScope, programSymbolTable, true)

		if err != nil {
			printExecError(err)
		}		
		
	}
}