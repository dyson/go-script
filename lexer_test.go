package main

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	input := `Replace(" ", "\n").Freq().First(1)`

	tests := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{FUNCTION, "Replace"},
		{LPAREN, "("},
		{STRING, " "},
		{COMMA, ","},
		{STRING, "\n"},
		{RPAREN, ")"},

		{PERIOD, "."},

		{FUNCTION, "Freq"},
		{LPAREN, "("},
		{RPAREN, ")"},

		{PERIOD, "."},

		{FUNCTION, "First"},
		{LPAREN, "("},
		{INT, "1"},
		{RPAREN, ")"},

		{EOF, "\000"},
	}

	l := newLexer(input)

	for i, tt := range tests {
		tok := l.getToken()

		if tok.ttype != tt.expectedType || tok.literal != tt.expectedLiteral {
			t.Logf("tests[%d] - tokentype expected=%q, got=%q", i, tt.expectedType, tok.ttype)
			t.Fatalf("tests[%d] - literal expected=%q, got=%q", i, tt.expectedLiteral, tok.literal)
		}
	}
}
