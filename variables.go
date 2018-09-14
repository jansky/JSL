package main

import (
	"fmt"
	"math/rand"
)

type symbolTableEntry struct {
	value langObject
	references int
}

type symbolTable struct {
	symbols map[uint64]*symbolTableEntry
}

func (s *symbolTable) insert(value langObject) (uint64, error) {

	key := rand.Uint64()

	for _, keyPresent := s.symbols[key]; keyPresent == true; _, keyPresent = s.symbols[key] {
		key = rand.Uint64()
	}

	newEntry := &symbolTableEntry{value, 1,}

	s.symbols[key] = newEntry

	return key, nil
}

func (s *symbolTable) retrieve(key uint64) (langObject, bool) {
	symbol, ok := s.symbols[key]

	if ok == false {
		return nil, false
	} else {
		return symbol.value, true
	}
}

func (s *symbolTable) retrieveOnPop(key uint64) (langObject, bool) {
	obj, ok := s.retrieve(key)

	if ok == false {
		return nil, false
	} else {
		s.decReference(key)
		return obj, true
	}
}

func (s *symbolTable) incReference(key uint64) error {
	symbol, ok := s.symbols[key]

	if ok == false {
		return fmt.Errorf("Bad symbol table reference: %X", key)
	}

	symbol.references++

	return nil
}

func (s *symbolTable) decReference(key uint64) error {
	symbol, ok := s.symbols[key]

	if ok == false {
		return fmt.Errorf("Bad symbol table reference: %X", key)
	}

	symbol.references--

	return nil
}

func (s *symbolTable) cleanUp() error {

	for k, v := range s.symbols {
		if v.references < 1 {
			delete(s.symbols, k)
		}
	}

	return nil

}

type langVariable struct {
	key uint64
	local bool
}

type variableScope struct {
	variables map[string]*langVariable
	parent *variableScope
}

func (v *variableScope) getWithLocal(name string, localsVisible bool) (uint64, bool) {

	variable, ok := v.variables[name]

	if ok == true {
		if variable.local && !localsVisible {
			ok = false
		}
	}

	if ok == false {
		if v.parent != nil {
			return v.parent.getWithLocal(name, false)
		} else {
			return 0, false
		}
	}

	return variable.key, true
}

func (v *variableScope) get(name string) (uint64, bool) {

	return v.getWithLocal(name, true)

}

func (v *variableScope) set(name string, key uint64) error {

	v.variables[name] = &langVariable{key,false,}

	return nil
}

func (v *variableScope) setLocal(name string, key uint64) error {
	v.variables[name] = &langVariable{key,true,}

	return nil
}

/*func (v *variableScope) setGlobal(name string, key uint64) error {
	
	parent := v.parent

	for ;parent != nil; parent = v.parent {
	}

	return parent.set(name, key)
}*/