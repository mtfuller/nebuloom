package parser

import (
	"unicode"

	"golang.org/x/exp/slices"
)

type Lexer struct {
	input []rune
	pos   int
}

type TokenType int32

const (
	TOKEN_EOF TokenType = iota
	TOKEN_UNKNOWN
	TOKEN_KEYWORD
	TOKEN_IDENTIFIER
	TOKEN_LCURLY
	TOKEN_RCURLY
	TOKEN_ASSIGN
	TOKEN_INTEGER
	TOKEN_STRING
)

var (
	tokenTypeMap = map[string]TokenType{
		"TOKEN_EOF":        TOKEN_EOF,
		"TOKEN_UNKNOWN":    TOKEN_UNKNOWN,
		"TOKEN_KEYWORD":    TOKEN_KEYWORD,
		"TOKEN_IDENTIFIER": TOKEN_IDENTIFIER,
		"TOKEN_LCURLY":     TOKEN_LCURLY,
		"TOKEN_RCURLY":     TOKEN_RCURLY,
		"TOKEN_ASSIGN":     TOKEN_ASSIGN,
		"TOKEN_INTEGER":    TOKEN_INTEGER,
		"TOKEN_STRING":     TOKEN_STRING,
	}

	tokenTypeToStr = []string{
		"TOKEN_EOF",
		"TOKEN_UNKNOWN",
		"TOKEN_KEYWORD",
		"TOKEN_IDENTIFIER",
		"TOKEN_LCURLY",
		"TOKEN_RCURLY",
		"TOKEN_ASSIGN",
		"TOKEN_INTEGER",
		"TOKEN_STRING",
	}
)

type Token struct {
	tokenType TokenType
	value     string
}

var whitespace = []rune{' ', '\n', '\t', '\r'}

func NewLexer(input string) Lexer {
	return Lexer{
		input: []rune(input),
		pos:   0,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		char := l.input[l.pos]
		if isWhitespace(char) {
			l.pos = l.pos + 1
		} else {
			return
		}
	}
}

func (l *Lexer) Lex() Token {
	lexeme := make([]rune, 0, 256)

	if l.pos >= len(l.input) {
		return Token{
			tokenType: TOKEN_EOF,
			value:     "EOF",
		}
	}

	l.skipWhitespace()

	isPossibleStringLiteral := false
	if l.input[l.pos] == '"' {
		isPossibleStringLiteral = true
		lexeme = append(lexeme, '"')
		l.pos = l.pos + 1
	}

	for l.pos < len(l.input) {
		char := l.input[l.pos]
		l.pos = l.pos + 1

		if char == '"' && isPossibleStringLiteral {
			lexeme = append(lexeme, char)
			return Token{
				tokenType: TOKEN_STRING,
				value:     string(lexeme),
			}
		}

		if !isWhitespace(char) || isPossibleStringLiteral {
			lexeme = append(lexeme, char)
		} else {
			return Token{
				tokenType: mapLexemeToTokenType(lexeme),
				value:     string(lexeme),
			}
		}
	}

	if len(lexeme) > 0 {
		return Token{
			tokenType: mapLexemeToTokenType(lexeme),
			value:     string(lexeme),
		}
	}

	return Token{
		tokenType: TOKEN_UNKNOWN,
		value:     "UNKNOWN",
	}
}

func isWhitespace(char rune) bool {
	return slices.Contains(whitespace, char)
}

func isEqual(a []rune, b []rune) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func mapLexemeToTokenType(lexeme []rune) TokenType {
	tokenType := TOKEN_UNKNOWN

	if isEqual(lexeme, []rune("->")) {
		return TOKEN_ASSIGN
	}

	isKeyword := true
	for _, char := range lexeme {
		if !unicode.IsLower(char) {
			isKeyword = false
		}
	}
	if isKeyword {
		return TOKEN_KEYWORD
	}

	for pos, char := range lexeme {
		if unicode.IsLetter(char) {
			if tokenType == TOKEN_UNKNOWN {
				tokenType = TOKEN_IDENTIFIER
			}
		} else if unicode.IsNumber(char) {
			if tokenType == TOKEN_UNKNOWN {
				tokenType = TOKEN_INTEGER
			}
		} else if char == '"' {
			if pos == 0 && lexeme[len(lexeme)-1] == '"' {
				tokenType = TOKEN_STRING
			}
		} else if char == '{' {
			return TOKEN_LCURLY
		} else if char == '}' {
			return TOKEN_RCURLY
		}
	}

	return tokenType
}
