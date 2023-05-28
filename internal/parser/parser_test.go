package parser

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type ParserFixture struct {
	Fixture struct {
		Tokens []struct {
			Token string
			Value string
		}
		Ast struct {
			Components map[string]map[string]interface{}
		}
	}
}

func TestParser(t *testing.T) {

	tests := []Test{
		{Name: "should correctly generate an AST when parsing component", FixtureFilepath: "./testdata/simple-component.yml"},
		{Name: "should correctly generate an AST when parsing store", FixtureFilepath: "./testdata/simple-store.yml"},
		{Name: "should correctly generate an AST when parsing service", FixtureFilepath: "./testdata/simple-service.yml"},
		{Name: "should correctly generate an AST when parsing function", FixtureFilepath: "./testdata/simple-function.yml"},
	}

	t.Run("should initialize a new Parser when NewParser() is called", func(t *testing.T) {

		// GIVEN a valid list of tokens
		tokens := []Token{
			{tokenType: TOKEN_EOF, value: "EOF"},
		}

		// WHEN NewParser() is called
		parser := NewParser(tokens)

		// THEN parser is initialized
		assert.NotNil(t, parser)

	})

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// GIVEN valid fixture file
			data, err := ioutil.ReadFile(test.FixtureFilepath)
			if err != nil {
				t.Fatalf("Failed to read fixture file: %s", err)
			}

			// AND valid test fixture
			var fixture ParserFixture
			err = yaml.Unmarshal(data, &fixture)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %s", err)
			}

			// AND valid token list
			tokenList := []Token{}
			for _, token := range fixture.Fixture.Tokens {
				tokenType, ok := tokenTypeMap[token.Token]
				assert.Truef(t, ok, "Expected %s to be a valid token", token.Token)

				tokenList = append(tokenList, Token{
					tokenType: tokenType,
					value:     token.Value,
				})
			}

			// AND valid parser
			parser := NewParser(tokenList)

			// WHEN parse() is called
			ast := parser.Parse()

			// THEN an AST is created
			assert.NotNil(t, ast)

			// AND the expected AST is generated
			assert.EqualValues(t, fixture.Fixture.Ast, ast)

		})
	}

}
