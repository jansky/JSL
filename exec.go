package main

import (
	"fmt"
)

func handleIdentifier(s *stack, v *variableScope, st *symbolTable, ident *langObjectIdentifier) error {
	stKey, ok := v.get(ident.name)

	switch ident.typ {
	case identifierDefault:
		if ok == false {
			return fmt.Errorf("Variable '%s' undefined in the local scope.", ident.name)
		}

		obj, stOk := st.retrieve(stKey)			

		if stOk == false {
			return fmt.Errorf("Unable to retrieve object with ID %X from the symbol table.", stKey)
		}

		/* Whenever we push a reference onto the stack, we must handle garbage collection */
		if obj.getType() == objectTypeReference {
			stErr := st.incReference(obj.(*langObjectReference).key)

			if stErr != nil {
				return fmt.Errorf("Unable to retrieve object with ID %X from the symbol table.", stKey)
			}
		}

		s.push(obj.copy())
	case identifierReference:

		if ok == true {
			stErr := st.incReference(stKey)

			if stErr != nil {
				return stErr
			}

			s.push(&langObjectReference{stKey,})
		} else {
			s.push(ident)
		}		
	}

	return nil
}

func (l *langObjectCodeBlock) exec(s *stack, v *variableScope, st *symbolTable, cleanUpLocal bool) error {

	for _, o := range l.code {

		/*fmt.Println(o.toString())
		fmt.Println("---")
		s.print()
		fmt.Println("---")*/

		switch {
		case o.getType() == objectTypeNumber|| o.getType() == objectTypeString || o.getType() == objectTypeBoolean:
			s.push(o)
		case o.getType() == objectTypeCodeBlock:
			codeBlock := o.(*langObjectCodeBlock)

			/* 	A code block has its parent scope set once, when it is first placed onto the stack.
				This allows recursive functions to work. */
			if codeBlock.parentScope == nil {
				codeBlock.parentScope = v
			}

			s.push(o)
		case o.getType() == objectTypeIdentifier:
			identErr := handleIdentifier(s, v, st, o.(*langObjectIdentifier))

			if identErr != nil {
				return identErr
			}
		case o.getType() == objectTypeOperation:
			err := performOperation(o.getValue().(operationType), s, v, st)

			if err != nil {
				return err
			}
		}

		/*s.print()
		fmt.Println("---\n")*/

	}

	if cleanUpLocal {
		for _, langVar := range v.variables {
			stErr := st.decReference(langVar.key)

			if stErr != nil {
				return stErr
			}
		}
	}

	st.cleanUp()

	return nil

} 