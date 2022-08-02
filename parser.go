package main

import "fmt"

func parse(input string) (*ast, error) {
	return newParser(newLexer(input)).parse()
}

type parser struct {
	l *lexer
	t token
}

func newParser(l *lexer) *parser {
	return &parser{l: l}
}

func (p *parser) nextToken() *parser {
	p.t = p.l.getToken()
	return p
}

func (p *parser) parse() (*ast, error) {
	tree := newAST()

	if p.nextToken().tokenIsType(EOF) {
		return tree, nil
	}

	for {
		f, err := p.parseFunction()
		if err != nil {
			return nil, err
		}

		tree.functions = append(tree.functions, *f)

		if p.nextToken().tokenIsType(EOF) {
			break
		}
		err = p.tokenMustType(PERIOD)
		if err != nil {
			return nil, err
		}
		p.nextToken()
	}

	return tree, nil
}

func (p *parser) tokenIsType(tt tokenType) bool {
	return p.t.ttype == tt
}

func (p *parser) tokenIsTypes(tts ...tokenType) bool {
	match := false
	for _, tt := range tts {
		if p.t.ttype == tt {
			match = true
			break
		}
	}

	return match
}

func (p *parser) tokenMustType(tt tokenType) error {
	if !p.tokenIsType(tt) {
		return fmt.Errorf("unexpected %s, expected %s", p.t.ttype, tt)
	}

	return nil
}

func (p *parser) tokenMustTypes(tts ...tokenType) error {
	if !p.tokenIsTypes(tts...) {
		return fmt.Errorf("unexpected %s, expected one of %s", p.t.ttype, tts)
	}

	return nil
}

func (p *parser) parseFunction() (*function, error) {
	err := p.tokenMustType(FUNCTION)
	if err != nil {
		return nil, err
	}

	name := p.t.literal

	err = p.nextToken().tokenMustType(LPAREN)
	if err != nil {
		return nil, err
	}

	args, err := p.nextToken().parseArguments()
	if err != nil {
		return nil, err
	}

	err = p.tokenMustType(RPAREN)
	if err != nil {
		return nil, err
	}

	f := function{
		name:      name,
		arguments: args,
	}

	return &f, nil
}

func (p *parser) parseArguments() ([]string, error) {
	args := []string{}

	if p.tokenIsTypes(RPAREN, EOF) {
		return args, nil
	}

	for {
		err := p.tokenMustTypes(STRING, INT)
		if err != nil {
			return nil, err
		}

		args = append(args, p.t.literal)

		if p.nextToken().tokenIsType(RPAREN) {
			break
		}
		err = p.tokenMustType(COMMA)
		if err != nil {
			return nil, err
		}
		p.nextToken()
	}

	return args, nil
}
