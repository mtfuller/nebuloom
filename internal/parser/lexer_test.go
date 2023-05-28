package parser

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type Test struct {
	Name            string
	FixtureFilepath string
}

type Fixture struct {
	Fixture struct {
		Source string
		Tokens []struct {
			Token string
			Value string
		}
	}
}

func TestLexer(t *testing.T) {

	tests := []Test{
		{Name: "should correctly extract lexemes when scanning component", FixtureFilepath: "./testdata/simple-component.yml"},
		{Name: "should correctly extract lexemes when scanning store", FixtureFilepath: "./testdata/simple-store.yml"},
		{Name: "should correctly extract lexemes when scanning service", FixtureFilepath: "./testdata/simple-service.yml"},
		{Name: "should correctly extract lexemes when scanning function", FixtureFilepath: "./testdata/simple-function.yml"},
	}

	t.Run("should initialize a new Lexer when NewLexer() is called", func(t *testing.T) {

		// GIVEN a valid Nebuloom program
		var prog = `
		component MyComponent {
			Field1 -> 123
		}`

		// WHEN NewLexer() is called
		var lexer = NewLexer(prog)

		// THEN lexer is initialized
		assert.NotNil(t, lexer)

		// AND lexer has program
		assert.Equal(t, lexer.input, []rune(prog))

	})

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// GIVEN valid fixture file
			data, err := ioutil.ReadFile(test.FixtureFilepath)
			if err != nil {
				t.Fatalf("Failed to read fixture file: %s", err)
			}

			// AND valid test fixture
			var fixture Fixture
			err = yaml.Unmarshal(data, &fixture)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %s", err)
			}

			// AND valid source code passed into the Lexer
			lexer := NewLexer(fixture.Fixture.Source)

			// THEN the Lexer scans and finds all expected lexemes
			for _, expectedLexeme := range fixture.Fixture.Tokens {
				actualLexeme := lexer.Lex()

				expectedTokenType, ok := tokenTypeMap[expectedLexeme.Token]

				assert.Truef(t, ok, "expected: %s, but got %s", expectedLexeme, actualLexeme)

				assert.Equal(t, expectedTokenType, actualLexeme.tokenType)
				assert.Equal(t, expectedLexeme.Value, actualLexeme.value)
			}

		})
	}

}
