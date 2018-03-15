package parse

import (
	"fmt"
	"strings"
)

// Query is the single query ast
type Query struct {
	Statement Statement
}

// Statement is the statement interface
type Statement interface {
	parse(*parser) error
}

// Field is the fields inside select
type Field struct {
	WildCard bool // aka '*'
	//Alias    string // add support for this? :))
	Table  string // the part before dot
	Column string // the column
}

// Fields is the collection of fields with order
type Fields []Field

// SelectStmt is the select query
type SelectStmt struct {
	Table  string
	Fields Fields

	Where Stack
}

// GetTokenString is a simple function to handle the quoted strings
func GetTokenString(t Item) string {
	v := t.Value()
	if t.Type() == ItemLiteral1 {
		v = strings.Trim(strings.Replace(t.Value(), `\'`, `'`, -1), "'")
	}
	if t.Type() == ItemLiteral2 {
		v = strings.Trim(strings.Replace(t.Value(), `\"`, `"`, -1), "\"")
	}
	return v
}

func (ss *SelectStmt) parseField(p *parser) (Field, error) {
	token := p.scanIgnoreWhiteSpace()
	if token.typ == ItemWildCard {
		return Field{WildCard: true}, nil
	}

	if token.typ == ItemAlpha || token.typ == ItemLiteral2 {
		ahead := p.scan() // white space is not allowed here
		if ahead.typ != ItemDot {
			p.reject()
			return Field{Column: GetTokenString(token)}, nil
		}
		ahead = p.scan()
		if ahead.typ == ItemAlpha || ahead.typ == ItemLiteral2 {
			return Field{
				Table:  GetTokenString(token),
				Column: GetTokenString(ahead),
			}, nil
		}
	}

	return Field{}, fmt.Errorf("unexpected token, %s", token)
}

func (ss *SelectStmt) parseFields(p *parser) error {
	for {
		field, err := ss.parseField(p)
		if err != nil {
			return err
		}
		ss.Fields = append(ss.Fields, field)

		comma := p.scanIgnoreWhiteSpace()
		if comma.typ != ItemComma {
			p.reject()
			break
		}
	}
	return nil
}

func (ss *SelectStmt) parse(p *parser) error {
	if err := ss.parseFields(p); err != nil {
		return err
	}

	t := p.scanIgnoreWhiteSpace()
	// must be from
	if t.typ != ItemFrom {
		return fmt.Errorf("unexpected %s , expected FROM or COMMA (,)", t)
	}

	t = p.scanIgnoreWhiteSpace()
	if t.typ != ItemAlpha && t.typ != ItemLiteral2 {
		return fmt.Errorf("unexpected input %s , need table name", t)
	}
	ss.Table = GetTokenString(t)

	if w := p.scanIgnoreWhiteSpace(); w.typ == ItemWhere {
		p.reject()
		var err error
		ss.Where, err = p.where()
		if err != nil {
			return err
		}
	}

	return nil
}

func newStatement(p *parser) (Statement, error) {
	start := p.scan()
	switch start.typ {
	case ItemSelect:
		sel := &SelectStmt{}
		err := sel.parse(p)
		if err != nil {
			return nil, err
		}
		return sel, nil
	default:
		return nil, fmt.Errorf("token %s is not a valid token", start.value)
	}
}