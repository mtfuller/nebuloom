package parser

import "strconv"

type Ast struct {
	Components map[string]map[string]interface{}
}

type Store struct {
	Schema map[string]string
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) nextToken() Token {
	nextPos := p.pos + 1
	if nextPos >= len(p.tokens) {
		return Token{
			tokenType: TOKEN_EOF,
			value:     "EOF",
		}
	} else {
		return p.tokens[nextPos]
	}
}

func (p *Parser) currentToken() Token {
	if p.pos >= len(p.tokens) {
		return Token{
			tokenType: TOKEN_EOF,
			value:     "EOF",
		}
	} else {
		return p.tokens[p.pos]
	}
}

func (p *Parser) expectToken(tokenType TokenType) Token {
	if p.pos >= len(p.tokens) {
		panic("Expected " + tokenTypeToStr[tokenType] + " type token, but none was found.")
	}

	nextToken := p.tokens[p.pos]
	p.pos = p.pos + 1

	if nextToken.tokenType != tokenType {
		panic("Expected " + tokenTypeToStr[tokenType] + " type token, but found " + tokenTypeToStr[nextToken.tokenType] + ".")
	}

	return nextToken
}

func (p *Parser) parseObject() interface{} {
	p.expectToken(TOKEN_LCURLY)

	currentToken := p.currentToken()
	fields := make(map[string]interface{})
	for currentToken.tokenType != TOKEN_RCURLY {
		field, value := p.parseFieldAssignment()
		fields[field] = value
		currentToken = p.currentToken()
	}

	p.expectToken(TOKEN_RCURLY)

	return fields
}

func (p *Parser) parseExpression() interface{} {
	currentToken := p.currentToken()
	nextToken := p.nextToken()

	if currentToken.tokenType == TOKEN_STRING {
		p.pos = p.pos + 1
		return currentToken.value
	} else if currentToken.tokenType == TOKEN_IDENTIFIER {
		p.pos = p.pos + 1
		if nextToken.tokenType == TOKEN_LCURLY {
			obj := p.parseObject()
			return map[string]interface{}{
				"type": "$METHOD",
				"ref":  currentToken.value,
				"args": obj,
			}
		} else {
			return currentToken.value
		}
	} else if currentToken.tokenType == TOKEN_INTEGER {
		p.pos = p.pos + 1
		i, err := strconv.Atoi(currentToken.value)
		if err != nil {
			panic(err)
		}
		return i
	} else if currentToken.tokenType == TOKEN_LCURLY {
		return p.parseObject()
	}

	return nil
}

func (p *Parser) parseFieldAssignment() (string, interface{}) {
	field := p.expectToken(TOKEN_IDENTIFIER)
	p.expectToken(TOKEN_ASSIGN)
	value := p.parseExpression()

	return field.value, value
}

func (p *Parser) parseComponent() (string, string, interface{}) {
	keyword := p.expectToken(TOKEN_KEYWORD)
	name := p.expectToken(TOKEN_IDENTIFIER)
	p.expectToken(TOKEN_LCURLY)

	currentToken := p.currentToken()
	var fields interface{}
	var components interface{}
	for currentToken.tokenType != TOKEN_RCURLY {
		if currentToken.tokenType == TOKEN_IDENTIFIER {
			if fields == nil {
				fields = make(map[string]interface{})
			}

			field, value := p.parseFieldAssignment()
			fields.(map[string]interface{})[field] = value
		} else if currentToken.tokenType == TOKEN_KEYWORD {
			if components == nil {
				components = make(map[string]interface{})
			}

			subkeyword, subname, subcomponent := p.parseComponent()

			_, ok := components.(map[string]interface{})[subkeyword]
			if !ok {
				components.(map[string]interface{})[subkeyword] = make(map[string]interface{})
			}

			subcomponents := components.(map[string]interface{})[subkeyword].(map[string]interface{})

			subcomponents[subname] = subcomponent
		}

		currentToken = p.currentToken()
	}

	p.expectToken(TOKEN_RCURLY)

	component := make(map[string]interface{})
	component["fields"] = fields
	component["components"] = components

	return keyword.value, name.value, component
}

func (p *Parser) Parse() interface{} {
	keyword, name, component := p.parseComponent()

	components := make(map[string]map[string]interface{})

	components[keyword] = make(map[string]interface{})

	components[keyword][name] = component

	return Ast{
		Components: components,
	}
}
