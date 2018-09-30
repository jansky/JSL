package main

import (
	"fmt"
)

func handleIdentifier(v *variableScope, st *symbolTable, ident *langObjectIdentifier) (langObject, error) {
	stKey, ok := v.get(ident.name)

	switch ident.typ {
	case identifierDefault:
		if ok == false {
			return nil, fmt.Errorf("Variable '%s' undefined in the local scope.", ident.name)
		}

		obj, stOk := st.retrieve(stKey)			

		if stOk == false {
			return nil, fmt.Errorf("Unable to retrieve object with ID %X from the symbol table.", stKey)
		}

		/* Whenever we push a reference onto the stack, we must handle garbage collection */
		if obj.getType() == objectTypeReference {
			stErr := st.incReference(obj.(*langObjectReference).key)

			if stErr != nil {
				return nil, fmt.Errorf("Unable to retrieve object with ID %X from the symbol table.", stKey)
			}
		}

		//s.push(obj.copy())
		return obj.copy(), nil
	case identifierReference:

		if ok == true {
			stErr := st.incReference(stKey)

			if stErr != nil {
				return nil, stErr
			}

			//s.push(&langObjectReference{stKey,})
			return &langObjectReference{stKey,}, nil
		} else {
			//s.push(ident)
			return ident, nil
		}		
	}

	return nil, nil
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
		case o.getType() == objectTypeList:
			if o.(*langObjectList).empty {
				s.push(o)
			} else {
				s.push(o.copy())
			}
		case o.getType() == objectTypeCodeBlock:
			codeBlock := o.(*langObjectCodeBlock)

			/* 	A code block has its parent scope set once, when it is first placed onto the stack.
				This allows recursive functions to work. */
			if codeBlock.parentScope == nil {
				codeBlock.parentScope = v
			}

			handleErr := o.(*langObjectCodeBlock).handleParentVariables(st)

			if handleErr != nil {
				return handleErr
			}

			s.push(o)
		case o.getType() == objectTypeIdentifier:
			identifier, identErr := handleIdentifier(v, st, o.(*langObjectIdentifier))

			if identErr != nil {
				return identErr
			}

			s.push(identifier)
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