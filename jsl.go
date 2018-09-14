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

		scanner.Scan()

		input := scanner.Text()

		if input == "exit" || input == "quit" {
			os.Exit(0)
		}
	
		_, items := lex("test", input)

		main, err := parseCodeBlock(&parser{items, 0,})

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}

		//main.print(0)

		execErr := main.exec(programStack, programVariableScope, programSymbolTable, false)

		if execErr != nil {
			fmt.Printf("Error: %s\n", execErr.Error())
			continue
		}

		/*stackObject, peekError := programStack.peek()

		if peekError == nil {
			fmt.Printf("%s\n", stackObject.toString())
		}*/

		programStack.print()
		
	}
}