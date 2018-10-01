package main

import (
	"fmt"
	"strconv"
	"regexp"
)

type parser struct {
	lexerItems chan item
	nestedLevel int
}

func parseString(str string) string {

	escape_runes := map[string]string{"n": "\n", "\"": "\"", "t": "\t", "/": "\\"}

	for k, v := range escape_runes {
		
		re := regexp.MustCompile("\\\\" + k)
		str = re.ReplaceAllString(str, v)

	}

	return str
}

func parseIdentifier(i item) langObject {
	
	if i.typ == itemIdentifierReference || i.typ == itemIdentifierReferenceAt || i.typ == itemIdentifierCall  || i.typ == itemIdentifierName {
		identifierTyp := identifierReference

		if i.typ == itemIdentifierReferenceAt {
			identifierTyp = identifierReferenceAt
		}

		if i.typ == itemIdentifierCall {
			identifierTyp = identifierCall
		}

		if i.typ == itemIdentifierName {
			identifierTyp = identifierName
		}

		return &langObjectIdentifier{identifierTyp, i.val,}
	} else {
		switch {
		case i.val == "clear":
			return &langObjectOperation{operationTypeClear,}
		case i.val == "asn":
			return &langObjectOperation{operationTypeAssign,}
		case i.val == "lasn":
			return &langObjectOperation{operationTypeLocalAssign,}
		case i.val == "true":
			return &langObjectBoolean{true,}
		case i.val == "false":
			return &langObjectBoolean{false,}
		case i.val == "dup":
			return &langObjectOperation{operationTypeDuplicate,}
		case i.val == "if":
			return &langObjectOperation{operationTypeIf,}
		case i.val == "for":
			return &langObjectOperation{operationTypeFor,}
		case i.val == "empty?":
			return &langObjectOperation{operationTypeListEmpty,}
		case i.val == "split":
			return &langObjectOperation{operationTypeListSplit,}
		case i.val == "pop":
			return &langObjectOperation{operationTypePop,}
		case i.val == "include":
			return &langObjectOperation{operationTypeInclude,}
		default:
			return &langObjectIdentifier{identifierDefault, i.val,}
		}
	}
}

func parseCodeBlock(p *parser) (*langObjectCodeBlock, error) {

	codeBlockItems := make([]langObject, 0)
	
	for i := range p.lexerItems {
		switch {
		case i.typ == itemError:
			return &langObjectCodeBlock{}, fmt.Errorf("Lexer error: %s", i.val)
		case i.typ == itemEOF:
			if p.nestedLevel == 0 {
				return &langObjectCodeBlock{codeBlockItems,nil}, nil
			} else {
				return &langObjectCodeBlock{}, fmt.Errorf("Unexpected end of file.")
			}
		case i.typ == itemOpenBlock:
			p.nestedLevel++
			codeBlock, err := parseCodeBlock(p)

			if err != nil {
				return &langObjectCodeBlock{}, err
			} else {
				codeBlockItems = append(codeBlockItems, codeBlock)
			}
			
		case i.typ == itemEndBlock:
			if p.nestedLevel < 1 {
				return &langObjectCodeBlock{}, fmt.Errorf("Unexpected end of block.")
			} else {
				p.nestedLevel--

				return &langObjectCodeBlock{codeBlockItems,nil}, nil
			}
		case i.typ == itemNumber:
			number, err := strconv.ParseFloat(i.val, 64)

			if err != nil {
				return &langObjectCodeBlock{}, fmt.Errorf("Error parsing number: '%s'.", i.val)
			} else {
				codeBlockItems = append(codeBlockItems, &langObjectNumber{number,})
			}
		case i.typ == itemString:
			codeBlockItems = append(codeBlockItems, &langObjectString{parseString(i.val),})
		case i.typ == itemPlus:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeAdd,})
		case i.typ == itemMinus:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeSubtract,})
		case i.typ == itemTimes:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeMultiply,})
		case i.typ == itemDividedBy:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeDivide,})
		case i.typ == itemExecute:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeExecute,})
		case i.typ == itemAt:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeAt,})
		case i.typ == itemNot:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeNot,})
		case i.typ == itemIdentifier || i.typ == itemIdentifierReference || i.typ == itemIdentifierReferenceAt || i.typ == itemIdentifierCall || i.typ == itemIdentifierName:
			codeBlockItems = append(codeBlockItems, parseIdentifier(i))
		case i.typ == itemCondition:
			switch i.val {
			case "=":
				codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeEquals,})
			case "<":
				codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeLess,})
			case ">":
				codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeGreater,})
			case "<=":
				codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeLessEquals,})
			case ">=":
				codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeGreaterEquals,})
			default:
				return &langObjectCodeBlock{}, fmt.Errorf("Unknown condition type '%s'.", i.val)
			}
		case i.typ == itemEmptyList:
			codeBlockItems = append(codeBlockItems, &langObjectList{true, nil, nil,})
		case i.typ == itemCons:
			codeBlockItems = append(codeBlockItems, &langObjectOperation{operationTypeCons,})
		default:
			return &langObjectCodeBlock{}, fmt.Errorf("Unknown lexer item type.")
		}
	}

	return &langObjectCodeBlock{codeBlockItems,nil}, nil	
}

