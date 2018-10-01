package main

import (
	"errors"
	"fmt"
	"os"
	"bufio"
)

func evaluateCondition(typ operationType, s *stack, v *variableScope, st *symbolTable) error {

	obj1, err1 := s.pop()

	if err1 != nil {
		return err1
	}

	obj2, err2 := s.pop()

	if err2 != nil {
		return err2
	}

	if obj1.getType() == objectTypeReference {
		err3 := st.decReference(obj1.getValue().(uint64))

		if err3 != nil {
			return err3
		}
	}

	if obj2.getType() == objectTypeReference {
		err4 := st.decReference(obj2.getValue().(uint64))

		if err4 != nil {
			return err4
		}
	}

	switch typ {
	case operationTypeEquals:
		condition, cErr := obj2.equals(obj1)

		if cErr != nil {
			return cErr
		}

		s.push(&langObjectBoolean{condition,})
	case operationTypeGreater:
		condition, cErr := obj2.greaterThan(obj1)

		if cErr != nil {
			return cErr
		}

		s.push(&langObjectBoolean{condition,})
	case operationTypeLess:
		condition, cErr := obj2.lessThan(obj1)

		if cErr != nil {
			return cErr
		}

		s.push(&langObjectBoolean{condition,})
	case operationTypeGreaterEquals:
		conditionGt, gtErr := obj2.greaterThan(obj1)

		if gtErr != nil {
			return gtErr
		}

		conditionEq, eqErr := obj2.equals(obj1)

		if eqErr != nil {
			return eqErr
		}

		s.push(&langObjectBoolean{conditionGt || conditionEq,})
	case operationTypeLessEquals:
		conditionLt, ltErr := obj2.lessThan(obj1)

		if ltErr != nil {
			return ltErr
		}

		conditionEq, eqErr := obj2.equals(obj1)

		if eqErr != nil {
			return eqErr
		}

		s.push(&langObjectBoolean{conditionLt || conditionEq,})
	default:
		return errors.New("Invalid condition type.")
	}

	return nil
}

func performAssign(typ operationType, s *stack, v *variableScope, st *symbolTable) error {

	reference, err1 := s.pop()

	if err1 != nil {
		return err1
	}

	val, err2 := s.pop()

	if err2 != nil {
		return err2
	}

	if reference.getType() == objectTypeIdentifier {
		ident := reference.(*langObjectIdentifier)

		if ident.typ != identifierReference {
			return errors.New("Expected, but did not receive an identifier reference.")
		}

		/* See if the reference points to an already defined variable */

		stKey, ok := v.get(ident.name)

		if ok == true {
			/* Garbage collection */
			err4 := st.decReference(stKey)

			if err4 != nil {
				return err4
			}
		}

		/* Insert the new variable value */
		valKey, err5 := st.insert(val)

		if err5 != nil {
			return err5
		}

		switch typ {
		case operationTypeAssign:
			v.set(ident.name, valKey)
		case operationTypeLocalAssign:
			v.setLocal(ident.name, valKey)
		}

		
	} else if reference.getType() == objectTypeReference {

		ref := reference.(*langObjectReference)

		err6 := st.decReference(ref.key)

		if err6 != nil {
			return err6
		}

		if typ == operationTypeLocalAssign {
			return errors.New("Cannot perform a local assign on a reference.")
		}

		/* Update the variable */

		st.symbols[ref.key].value = val

	} else {
		return errors.New("Expected, but did not receive an identifier reference or reference.")
	}

	return nil
}

func performOperation(typ operationType, s *stack, v *variableScope, st *symbolTable) error {

	/*fmt.Println("---")
	s.print()
	fmt.Println("---")*/

	switch typ {
	case operationTypeAdd:
		num1, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		num2, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		var resultObject langObject

		if num1.getType() == objectTypeNumber && num2.getType() == objectTypeNumber {

			result := num1.getValue().(float64) + num2.getValue().(float64)

			resultObject = &langObjectNumber{result,}

		} else {
			result := num2.toString() + num1.toString()

			resultObject = &langObjectString{result,}
		}

		err3 := s.push(resultObject)

		if err3 != nil {
			return err3
		}		

		
	case operationTypeSubtract:
		num1, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if num1.getType() != objectTypeNumber {
			return errors.New("Expected first item to be number.")
		}

		num2, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		if num2.getType() != objectTypeNumber {
			return errors.New("Expected second item to be number.")
		}

		result := num2.getValue().(float64) - num1.getValue().(float64)

		err3 := s.push(&langObjectNumber{result,})

		if err3 != nil {
			return err3
		}
	case operationTypeMultiply:
		num1, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if num1.getType() != objectTypeNumber {
			return errors.New("Expected first item to be number.")
		}

		num2, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		if num2.getType() != objectTypeNumber {
			return errors.New("Expected second item to be number.")
		}

		result := num1.getValue().(float64) * num2.getValue().(float64)

		err3 := s.push(&langObjectNumber{result,})

		if err3 != nil {
			return err3
		}
	case operationTypeDivide:
		num1, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if num1.getType() != objectTypeNumber {
			return errors.New("Expected first item to be number.")
		}

		num2, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		if num2.getType() != objectTypeNumber {
			return errors.New("Expected second item to be number.")
		}

		result := num2.getValue().(float64) / num1.getValue().(float64)

		err3 := s.push(&langObjectNumber{result,})

		if err3 != nil {
			return err3
		}
	case operationTypeExecute:
		block, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if block.getType() != objectTypeCodeBlock {
			return errors.New("Expected, but did not receive a code block.")
		}

		newScope := &variableScope{make(map[string]*langVariable),block.(*langObjectCodeBlock).parentScope,}
		err2 := block.(*langObjectCodeBlock).exec(s, newScope, st, true)

		if err2 != nil {
			return err2
		}
	case operationTypeClear:
		s.clear()
		return nil
	case operationTypeAssign, operationTypeLocalAssign:
		return performAssign(typ, s, v, st)		
	case operationTypeAt:
		/*identObj, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if identObj.getType() != objectTypeIdentifier {
			return errors.New("Expected, but did not receive an identifier reference.")
		}

		ident := identObj.(*langObjectIdentifier)

		if ident.typ != identifierReference {
			return fmt.Errorf("Cannot dereference an object that it is not an identifier reference.")
		}*/

		refObj, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if refObj.getType() != objectTypeReference {
			return errors.New("Expected, but did not receive a reference.")
		}

		ref := refObj.(*langObjectReference)

		/* Garbage Collection */

		/*stKey, ok := v.get(ident.name)

		if ok == false {
			return fmt.Errorf("Identifier '%s' does not reference any object.", ident.name)
		}*/

		err2 := st.decReference(ref.key)

		if err2 != nil {
			return err2
		}

		deRefObj, deRefOk := st.retrieve(ref.key)

		if deRefOk != true {
			return fmt.Errorf("The symbol table does not contain an entry for: %X.", ref.key)
		}

		//return handleIdentifier(s, v, st, deRefIdent.(*langObjectIdentifier))

		if deRefObj.getType() == objectTypeReference {
			stErr := st.incReference(deRefObj.(*langObjectReference).key)

			if stErr != nil {
				return fmt.Errorf("Unable to retrieve object with ID %X from the symbol table.", deRefObj.(*langObjectReference).key)
			}
		}

		s.push(deRefObj)
	case operationTypeNot:
		boolObj, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if boolObj.getType() != objectTypeBoolean {
			return errors.New("Expected, but did not receive a boolean.")
		}

		s.push(&langObjectBoolean{!(boolObj.(*langObjectBoolean).val)})
	case operationTypeEquals, operationTypeGreater, operationTypeLess, operationTypeGreaterEquals, operationTypeLessEquals:
		err1 := evaluateCondition(typ, s, v, st)

		if err1 != nil {
			return err1
		}
	case operationTypeDuplicate:
		obj, err1 := s.peek()

		if err1 != nil {
			return err1
		}

		if obj.getType() == objectTypeReference {
			stErr := st.incReference(obj.(*langObjectReference).key)

			if stErr != nil {
				return stErr
			}
		}

		s.push(obj.copy())
	case operationTypeIf:
		boolObj, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		codeBlockObj, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		if boolObj.getType() != objectTypeBoolean {
			return errors.New("Expected, but did not receive a boolean.")
		}

		if codeBlockObj.getType() != objectTypeCodeBlock {
			return errors.New("Expected, but did not recieve a code block.")
		}

		if boolObj.getValue().(bool) {

			newScope := &variableScope{make(map[string]*langVariable),codeBlockObj.(*langObjectCodeBlock).parentScope,}
			err3 := codeBlockObj.(*langObjectCodeBlock).exec(s, newScope, st, true)

			if err3 != nil {
				return err3
			}

		}
	case operationTypeFor:
		afterCodeBlock, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		bodyCodeBlock, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		conditionCodeBlock, err3 := s.pop()

		if err3 != nil {
			return err3
		}

		initialCodeBlock, err4 := s.pop()

		if err4 != nil {
			return err4
		}

		if afterCodeBlock.getType() != objectTypeCodeBlock || bodyCodeBlock.getType() != objectTypeCodeBlock || conditionCodeBlock.getType() != objectTypeCodeBlock || initialCodeBlock.getType() != objectTypeCodeBlock {
			return errors.New("Expected, but did not receive 4 code blocks.")
		}

		/*	Now we set up the variable scopes so that the for loop works properly:

			bodyCodeBlock.parentScope
				conditionScope (shared by afterCodeBlock, conditionCodeBlock, and initialCodeBlock)
					bodyCodeBlock
		*/

		conditionScope := &variableScope{make(map[string]*langVariable), bodyCodeBlock.(*langObjectCodeBlock).parentScope,}

		initErr := initialCodeBlock.(*langObjectCodeBlock).exec(s, conditionScope, st, false)

		if initErr != nil {
			return initErr
		}

		for true {
			
			condErr := conditionCodeBlock.(*langObjectCodeBlock).exec(s, conditionScope, st, false)

			if condErr != nil {
				return condErr
			}

			boolObj, boolErr := s.pop()

			if boolErr != nil {
				return boolErr
			}

			if boolObj.getType() != objectTypeBoolean {
				return errors.New("Expected, but did not get a boolean.")
			}

			if boolObj.(*langObjectBoolean).val == false {
				break
			}

			bodyScope := &variableScope{make(map[string]*langVariable), conditionScope,}

			bodyErr := bodyCodeBlock.(*langObjectCodeBlock).exec(s, bodyScope, st, true)

			if bodyErr != nil {
				return bodyErr
			}

			afterErr := afterCodeBlock.(*langObjectCodeBlock).exec(s, conditionScope, st, false)

			if afterErr != nil {
				return afterErr
			}

		}

		/* We must clean up the condition scope manually */

		for _, langVar := range conditionScope.variables {
			stErr := st.decReference(langVar.key)

			if stErr != nil {
				return stErr
			}
		}

		st.cleanUp()
	case operationTypeListEmpty:
		list, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if list.getType() != objectTypeList {
			return errors.New("Expected, but did not get a list.")
		}

		s.push(&langObjectBoolean{list.(*langObjectList).empty,})
	case operationTypeCons:

		value, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		list, err2 := s.pop()

		if err2 != nil {
			return err2
		}

		if list.getType() != objectTypeList {
			return errors.New("Expected, but did not get an object and a list.")
		}

		s.push(&langObjectList{
			false,
			value,
			list.copy().(*langObjectList),
		})
	case operationTypeListSplit:

		list, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if list.getType() != objectTypeList {
			return errors.New("Expected, but did not get a list.")
		}

		if list.(*langObjectList).empty {
			return errors.New("Unable to split an empty list.")
		}

		s.push(list.(*langObjectList).tail.copy())
		s.push(list.(*langObjectList).head)
	case operationTypePop:

		_, err1 := s.pop()

		if err1 != nil {
			return err1
		}
	case operationTypeInclude:

		filePath, err1 := s.pop()

		if err1 != nil {
			return err1
		}

		if filePath.getType() != objectTypeString {
			return errors.New("Expected, but did not get a string.")
		}

		includeFile, err2 := os.Open(filePath.(*langObjectString).val)

		if err2 != nil {
			return err2
		}

		fileScanner := bufio.NewScanner(includeFile)

		includeFileContents := ""

		for fileScanner.Scan() {
			includeFileContents += fileScanner.Text() + " "
		}

		err3 := evalString(includeFileContents, s, v, st, false)

		if err3 != nil {
			return err3
		}

		includeFile.Close()
	}

	

	return nil
}