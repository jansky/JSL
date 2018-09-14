package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

const eof = -1
const identifierRunes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"

const (
	itemError itemType = iota
	itemEOF
	itemIdentifier
	itemIdentifierReference
	itemIdentifierReferenceAt
	itemIdentifierCall
	itemOpenBlock
	itemEndBlock
	itemNumber
	itemString
	itemPlus
	itemMinus
	itemTimes
	itemDividedBy
	itemExecute
	itemAt
	itemNot
	itemCondition
)

type item struct {
	typ itemType
	val string
}

func (item *item) print() {
	itemTypeString := "unknown"
	printValue := true

	switch item.typ {
	case itemError:
		itemTypeString = "itemError"
	case itemEOF:
		itemTypeString = "EOF"
		printValue = false
	case itemIdentifier:
		itemTypeString = "itemIdentifier"
	case itemIdentifierReference:
		itemTypeString = "itemIdentifierReference"
	case itemIdentifierReferenceAt:
		itemTypeString = "itemIdentifierReferenceAt"
	case itemIdentifierCall:
		itemTypeString = "itemIdentifierCall"
	case itemOpenBlock:
		itemTypeString = "itemOpenBlock"
	case itemEndBlock:
		itemTypeString = "itemEndBlock"
	case itemNumber:
		itemTypeString = "itemNumber"
	case itemString:
		itemTypeString = "itemString"
	case itemPlus:
		itemTypeString = "itemPlus"
	case itemMinus:
		itemTypeString = "itemMinus"
	case itemTimes:
		itemTypeString = "itemTimes"
	case itemDividedBy:
		itemTypeString = "itemDividedBy"
	}

	if printValue {
		fmt.Printf("%s: %s\n", itemTypeString, item.val)
	} else 	{
		fmt.Printf("%s\n", itemTypeString)
	}
}

type lexer struct {
	name string
	input string
	start int
	pos int
	width int
	items chan item
}

type stateFn func(*lexer) stateFn

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (rune) {

	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width

	l.pos += l.width

	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item {
		itemError,
		fmt.Sprintf(format, args...),
	}

	return nil
}

func (l *lexer) run() {

	for state := lexCode; state != nil; {
		state = state(l)
	}
	
	close(l.items)
}

func lexIdentifier(l *lexer) stateFn {

	identifierType := itemIdentifier

	switch prefix := l.next(); {
	case prefix == '\'':
		identifierType = itemIdentifierReference
	case prefix == '!':
		identifierType = itemIdentifierCall
	case prefix == '@':
		identifierType = itemIdentifierReferenceAt
	}

	if identifierType == itemIdentifier {
		l.backup()
	} else {
		l.ignore() // We don't want the ' or the ! in the identifier name
	}

	l.acceptRun(identifierRunes)

	if l.start == l.pos {
		return l.errorf("Empty identifier at position %d", l.start)
	}
	
	l.emit(identifierType)
	return lexCode
}

func lexNumber(l *lexer) stateFn {

	l.accept("-")

	digits := "0123456789"

	if l.accept("0") && l.accept("xX") {
		digits = "0123456789ABCDEFabcdef"
	}

	l.acceptRun(digits)

	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(itemNumber)
	return lexCode
}

func lexQuotedString(l *lexer) stateFn {

	escape_sequence := false

	for c := l.next(); c != '"' || escape_sequence; c = l.next() {
		if c == '\\' {
			escape_sequence = true
		} else {
			escape_sequence = false
		}
	}

	l.backup()
	l.emit(itemString)

	l.next()
	l.ignore() // We don't want the end quote

	return lexCode
}

func lexCode(l *lexer) stateFn {

	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case unicode.IsSpace(r):
		l.ignore()
		return lexCode
	case r == '+':
		l.emit(itemPlus)
		return lexCode
	case r == '*':
		l.emit(itemTimes)
		return lexCode
	case r == '/':
		l.emit(itemDividedBy)
		return lexCode
	case r == '-':
		if n := l.peek(); '0' <= n && n <= '9' {
			l.backup()
			return lexNumber
		}

		l.emit(itemMinus)
		return lexCode
	case r == '{':
		l.emit(itemOpenBlock)
		return lexCode
	case r == '}':
		l.emit(itemEndBlock)
		return lexCode
	case r == '@':
		l.emit(itemAt)
		return lexCode
	case r == '~':
		l.emit(itemNot)
		return lexCode
	case r == '=':
		l.emit(itemCondition)
		return lexCode
	case r == '<' || r == '>':
		if l.next() != '=' {
			l.backup()
		}

		l.emit(itemCondition)
		return lexCode
	case r == '"':
		l.ignore() // We don't want the beginning quote
		return lexQuotedString
	case '0' <= r && r <= '9':
		l.backup()
		return lexNumber
	case strings.IndexRune(identifierRunes, r) >= 0:
		l.backup()
		return lexIdentifier
	case r == '\'':
		l.backup()
		return lexIdentifier
	case r == '!':
		if i := l.peek(); strings.IndexRune(identifierRunes, i) >= 0 {
			l.backup()
			return lexIdentifier
		}

		l.emit(itemExecute)
		return lexCode
	default:
		return l.errorf("Unexpected %q at position %d", r, l.start)
	}
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer {
		name: name,
		input: input,
		items: make(chan item),
	}

	go l.run()
	return l, l.items
}

