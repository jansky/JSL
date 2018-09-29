package main

import (
	"fmt"
	"errors"
)

type langObjectType int

const (
	objectTypeString langObjectType = iota
	objectTypeNumber
	objectTypeBoolean
	objectTypeOperation
	objectTypeCodeBlock
	objectTypeIdentifier
	objectTypeReference
	objectTypeList
	objectTypeError
)

type langObject interface {
	getType() langObjectType
	getValue() interface{}
	setValue(interface{})
	toString() string
	copy() langObject
	equals(langObject) (bool, error)
	greaterThan(langObject) (bool, error)
	lessThan(langObject) (bool, error)
}

type stack struct {
	contents []langObject
}

func (s *stack) push(l langObject) error {

	s.contents = append(s.contents, l)
	return nil
}

func (s *stack) pop() (langObject, error) {
	var i langObject

	if len(s.contents) < 1 {
		return i, errors.New("Stack underflow")
	}

	

	i, s.contents = s.contents[len(s.contents) - 1], s.contents[:len(s.contents)-1]

	return i, nil
}

func (s *stack) peek() (langObject, error) {
	i, e := s.pop()

	if e != nil {
		return nil, e
	}

	s.push(i)

	return i, nil
}

func (s *stack) clear() error {
	s.contents = make([]langObject, 0)

	return nil
}

func (s *stack) print() {
	for i := len(s.contents) - 1; i >= 0; i-- {
		fmt.Printf("%s\n", s.contents[i].toString())
	}
}

/* Strings */

type langObjectString struct {
	val string
}

func (l *langObjectString) getType() langObjectType {
	return objectTypeString
}

func (l *langObjectString) getValue() interface{} {
	return l.val
}

func (l *langObjectString) setValue(val interface{}) {
	l.val = val.(string)
}

func (l *langObjectString) toString() string {
	return l.val
}

func (l * langObjectString) copy() langObject {
	newVal := l.val

	return &langObjectString{newVal,}
}

func (l *langObjectString) equals(o langObject) (bool, error) {

	switch o.getType() {
	case objectTypeString:
		return (l.val == o.(*langObjectString).val), nil
	default:
		return false, nil
	}

}

func (l *langObjectString) greaterThan(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeString:
		return (l.val > o.(*langObjectString).val), nil
	default:
		return false, nil
	}
}

func (l *langObjectString) lessThan(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeString:
		return (l.val < o.(*langObjectString).val), nil
	default:
		return false, nil
	}
}

/* Numbers */

type langObjectNumber struct {
	val float64
}

func (l * langObjectNumber) getType() langObjectType {
	return objectTypeNumber
}

func (l *langObjectNumber) getValue() interface{} {
	return l.val
}

func (l *langObjectNumber) setValue(val interface{}) {
	l.val = val.(float64)
}

func (l *langObjectNumber) toString() string {
	return fmt.Sprintf("%f", l.val)
}

func (l *langObjectNumber) copy() langObject {
	return &langObjectNumber{l.val,}
}

func (l *langObjectNumber) equals(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeNumber:
		return (l.val == o.(*langObjectNumber).val), nil
	default:
		return false, nil
	}
}

func (l *langObjectNumber) greaterThan(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeNumber:
		return (l.val > o.(*langObjectNumber).val), nil
	default:
		return false, nil
	}
}

func (l *langObjectNumber) lessThan(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeNumber:
		return (l.val < o.(*langObjectNumber).val), nil
	default:
		return false, nil
	}
}

/* Boolean */

type langObjectBoolean struct {
	val bool
}

func (l *langObjectBoolean) getType() langObjectType {
	return objectTypeBoolean
}

func (l *langObjectBoolean) getValue() interface{} {
	return l.val
}

func (l *langObjectBoolean) setValue(val interface{}) {
	l.val = val.(bool)
}

func (l *langObjectBoolean) toString() string {
	if l.val {
		return "true"
	} else {
		return "false"
	}
}

func (l *langObjectBoolean) copy() langObject {
	return &langObjectBoolean{l.val,}
}

func (l *langObjectBoolean) equals(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeBoolean:
		return (l.val == o.(*langObjectBoolean).val), nil
	default:
		return false, nil
	}
}

func (l *langObjectBoolean) lessThan(o langObject) (bool, error) {
	return false, errors.New("Operation less than cannot be applied to boolean.")
}

func (l *langObjectBoolean) greaterThan(o langObject) (bool, error) {
	return false, errors.New("Operation greater than cannot be applied to boolean.")
}

/* Operation */

type operationType int

const (
	operationTypeAdd operationType = iota
	operationTypeSubtract
	operationTypeMultiply
	operationTypeDivide
	operationTypeExecute
	operationTypeClear
	operationTypeAssign
	operationTypeLocalAssign
	operationTypeAt
	operationTypeNot
	operationTypeEquals
	operationTypeGreater
	operationTypeLess
	operationTypeGreaterEquals
	operationTypeLessEquals
	operationTypeIf
	operationTypeFor
	operationTypeDuplicate
	operationTypeCons
	operationTypeListEmpty
	operationTypeListSplit
	operationTypePop
)

type langObjectOperation struct {
	val operationType
}

func (l *langObjectOperation) getType() langObjectType {
	return objectTypeOperation
}

func (l *langObjectOperation) getValue() interface{} {
	return l.val
}

func (l *langObjectOperation) setValue(val interface{}) {
	l.val = val.(operationType)
}

func(l *langObjectOperation) toString() string {
	operationName := "unknown"

	switch l.val {
	case operationTypeAdd:
		operationName = "add"
	case operationTypeSubtract:
		operationName = "subtract"
	case operationTypeMultiply:
		operationName = "multiply"
	case operationTypeDivide:
		operationName = "divide"
	case operationTypeExecute:
		operationName = "execute"
	case operationTypeClear:
		operationName = "clear"
	case operationTypeAssign:
		operationName = "assign"
	case operationTypeLocalAssign:
		operationName = "local assign"
	case operationTypeAt:
		operationName = "at"
	case operationTypeNot:
		operationName = "not"
	case operationTypeEquals:
		operationName = "equals"
	case operationTypeLess:
		operationName = "less than"
	case operationTypeGreater:
		operationName = "greater than"
	case operationTypeGreaterEquals:
		operationName = "greater than or equals"
	case operationTypeLessEquals:
		operationName = "less than or equals"
	case operationTypeIf:
		operationName = "if"
	case operationTypeFor:
		operationName = "for"
	case operationTypeDuplicate:
		operationName = "duplicate"
	case operationTypeCons:
		operationName = "cons"
	case operationTypeListEmpty:
		operationName = "empty?"
	case operationTypeListSplit:
		operationName = "split"
	case operationTypePop:
		operationName = "pop"
	}

	return fmt.Sprintf("<Operation: %s>", operationName)
}

func (l *langObjectOperation) copy() langObject {
	return &langObjectOperation{l.val,}
}

func (l *langObjectOperation) equals(o langObject) (bool, error){
	switch o.getType() {
	case objectTypeOperation:
		return (l.val == o.(*langObjectOperation).val), nil
	default:
		return false, nil
	}
}

func (l *langObjectOperation) greaterThan(o langObject) (bool, error) {
	return false, errors.New("Operation greater than cannot be applied to operation.")
}

func (l *langObjectOperation) lessThan(o langObject) (bool, error) {
	return false, errors.New("Operation less than cannot be applied to operation.")
}

/* Code Block */

type langObjectCodeBlock struct {
	code []langObject
	parentScope *variableScope
}

func (l *langObjectCodeBlock) getType() langObjectType {
	return objectTypeCodeBlock
}

func (l *langObjectCodeBlock) getValue() interface{} {
	return l.code
}

func (l *langObjectCodeBlock) setValue(code interface{}) {
	l.code = code.([]langObject)
}

func (l * langObjectCodeBlock) toString() string {
	return "<CodeBlock>"
}

func (l *langObjectCodeBlock) print(tabs int) {
	for _, obj := range l.getValue().([]langObject) {

		for i := 0; i < tabs; i++ {
			fmt.Printf("\t")
		}
		
		fmt.Printf("%s\n", obj.toString())

		if obj.getType() == objectTypeCodeBlock {
			obj.(*langObjectCodeBlock).print(tabs + 1)
		}

	}
}

func (l *langObjectCodeBlock) copy() langObject {
	newCode := make([]langObject,0)

	for _, obj := range l.code {
		newCode = append(newCode, obj.copy())
	}

	return &langObjectCodeBlock{newCode,l.parentScope,}
}

func (l *langObjectCodeBlock) equals(o langObject) (bool, error) {
	return false, errors.New("Code block objects are not comparable.")
}

func (l *langObjectCodeBlock) greaterThan(o langObject) (bool, error) {
	return false, errors.New("Code block objects are not comparable.")
}

func (l *langObjectCodeBlock) lessThan(o langObject) (bool, error) {
	return false, errors.New("Code block objects are not comparable.")
}

/* Identifier */

type identifierType int

const (
	identifierDefault identifierType = iota
	identifierReference
	identifierReferenceAt
	identifierCall
)

type langObjectIdentifier struct {
	typ identifierType
	name string
}

func (l *langObjectIdentifier) getType() langObjectType {
	return objectTypeIdentifier
}

func (l *langObjectIdentifier) getValue() interface{} {
	return l
}

func (l *langObjectIdentifier) setValue(s interface{}) {
	l = s.(*langObjectIdentifier)
}

func (l *langObjectIdentifier) toString() string {
	prefix := ""

	switch l.typ {
	case identifierReference:
		prefix = "'"
	case identifierReferenceAt:
		prefix = "@"
	case identifierCall:
		prefix = "!"
	}

	return prefix + l.name
}

func (l *langObjectIdentifier) copy() langObject {
	newName := l.name

	return &langObjectIdentifier{l.typ, newName,}
}

func (l *langObjectIdentifier) equals(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeIdentifier:
		return ((l.typ == o.(*langObjectIdentifier).typ) && (l.name == o.(*langObjectIdentifier).name)), nil
	default:
		return false, nil
	}
}

func (l *langObjectIdentifier) greaterThan(o langObject) (bool, error) {
	return false, errors.New("Operation greater than cannot be applied to identifier.")
}

func (l *langObjectIdentifier) lessThan(o langObject) (bool, error) {
	return false, errors.New("Operation less than cannot be applied to identifier.")
}

/* Reference */

type langObjectReference struct {
	key uint64
}

func (l *langObjectReference) getType() langObjectType {
	return objectTypeReference
}

func (l *langObjectReference) getValue() interface{} {
	return l.key
}

func (l *langObjectReference) setValue(key interface{}) {
	l.key = key.(uint64)
}

func (l *langObjectReference) toString() string {
	return fmt.Sprintf("<Reference: %X>", l.key)
}

func (l *langObjectReference) copy() langObject {
	return &langObjectReference{l.key,}
}

func (l *langObjectReference) equals(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeReference:
		return (l.key == o.(*langObjectReference).key), nil
	default:
		return false, nil
	}
}

func (l *langObjectReference) greaterThan(o langObject) (bool, error) {
	return false, errors.New("Operation greater than cannot be applied to reference.")
}

func (l *langObjectReference) lessThan(o langObject) (bool, error) {
	return false, errors.New("Operation less than cannot be applied to reference.")
}

/* Error */

type langObjectError struct {
	message string
}

func (l *langObjectError) getType() langObjectType {
	return objectTypeError
}

func (l *langObjectError) getValue() interface{} {
	return l.message
}

func (l *langObjectError) setValue(message interface{}) {
	l.message = message.(string)
}

func (l *langObjectError) toString() string {
	return l.message
}

func (l *langObjectError) copy() langObject {
	newMessage := l.message

	return &langObjectError{newMessage,}
}

func (l *langObjectError) equals(o langObject) (bool, error) {
	switch o.getType() {
	case objectTypeError:
		return (l.message == o.(*langObjectError).message), nil
	default:
		return false, nil
	}
}

func (l *langObjectError) greaterThan(o langObject) (bool, error) {
	return false, errors.New("Operation greater than cannot be applied to error.")
}

func (l *langObjectError) lessThan(o langObject) (bool, error) {
	return false, errors.New("Operation less than cannot be applied to error.")
}

/* Lists */

type langObjectList struct {
	empty bool
	head langObject
	tail *langObjectList
}

func (l *langObjectList) getType() langObjectType {
	return objectTypeList
}

func (l *langObjectList) getValue() interface{} {
	return l
}

func (l *langObjectList) setValue(list interface{}) {
	l = list.(*langObjectList)
}

func (l *langObjectList) toString() string {
	if l.empty {
		return "<>"
	} else {
		return "<...>"
	}
}

func (l *langObjectList) copy() langObject {

	var head_copied langObject
	tail_copied := &langObjectList{}

	if l.tail != nil {
		tail_copied = l.tail.copy().(*langObjectList)
	}

	if l.head != nil {
		head_copied = l.head.copy()
	}

	return &langObjectList{
		l.empty,
		head_copied,
		tail_copied,
	}

}

func (l *langObjectList) equals(obj langObject) (bool, error) {
	return false, errors.New("Operation equals cannot be applied to list.")
}

func (l *langObjectList) greaterThan(obj langObject) (bool, error) {
	return false, errors.New("Operation greater than cannot be applied to list.")
}

func (l *langObjectList) lessThan(obj langObject) (bool, error) {
	return false, errors.New("Operation less than cannot be applied to list.")
}

