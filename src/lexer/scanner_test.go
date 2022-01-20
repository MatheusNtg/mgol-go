package lexer

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanNumToken(t *testing.T) {
	testCases := []struct {
		name           string
		preparedText   string
		expectedTokens []Token
	}{
		{
			name:           "Integer number",
			preparedText:   "1",
			expectedTokens: []Token{NewToken(NUM, "1", INTEGER)},
		},
		{
			name:           "Real number",
			preparedText:   "1.0",
			expectedTokens: []Token{NewToken(NUM, "1.0", REAL)},
		},
		{
			name:           "Integer with capital exponential positive",
			preparedText:   "1E+0",
			expectedTokens: []Token{NewToken(NUM, "1E+0", INTEGER)},
		},
		{
			name:           "Integer with lower exponential positive",
			preparedText:   "1e+0",
			expectedTokens: []Token{NewToken(NUM, "1e+0", INTEGER)},
		},
		{
			name:           "Integer with capital exponential negative",
			preparedText:   "1E-0",
			expectedTokens: []Token{NewToken(NUM, "1E-0", INTEGER)},
		},
		{
			name:           "Integer with lower exponential negative",
			preparedText:   "1e-0",
			expectedTokens: []Token{NewToken(NUM, "1e-0", INTEGER)},
		},
		{
			name:         "Incomplete real number with capital exponential positive",
			preparedText: "1.E+0",
			expectedTokens: []Token{
				ERROR_TOKEN,
				NewToken(IDENTIFIER, "E", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(NUM, "0", INTEGER),
			},
		},
		// {
		// 	name:         "Incomplete real number with capital exponential positive",
		// 	preparedText: "1.e+0",
		// 	expectedTokens: []Token{
		// 		NewToken(NUM, "1", INTEGER),
		// 		ERROR_TOKEN,
		// 	},
		// },
		// {
		// 	name:         "Incomplete real number with capital exponential negative",
		// 	preparedText: "1.E-0",
		// 	expectedTokens: []Token{
		// 		NewToken(NUM, "1", INTEGER),
		// 		ERROR_TOKEN,
		// 	},
		// },
		// {
		// 	name:         "Incomplete real number with lower exponential negative",
		// 	preparedText: "1.e-0",
		// 	expectedTokens: []Token{
		// 		NewToken(NUM, "1", INTEGER),
		// 		ERROR_TOKEN,
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())
			tokens := []Token{}
			for {
				token := scanner.Scan()
				if token == EOF_TOKEN {
					break
				}
				tokens = append(tokens, token)
			}

			require.Equal(t, tc.expectedTokens, tokens)
		})
	}
}

func TestScanIdToken(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken Token
	}{
		{
			name:          "Identifier with number",
			preparedText:  "id1",
			expectedToken: NewToken(IDENTIFIER, "id1", NULL),
		},
		{
			name:          "Identifier with underline and number",
			preparedText:  "id_1",
			expectedToken: NewToken(IDENTIFIER, "id_1", NULL),
		},
		{
			name:          "Identifier with underline and more than one number",
			preparedText:  "id_123",
			expectedToken: NewToken(IDENTIFIER, "id_123", NULL),
		},
		{
			name:          "Identifier with underline and more than one number and more than one character",
			preparedText:  "id_123_id",
			expectedToken: NewToken(IDENTIFIER, "id_123_id", NULL),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())
			token := scanner.Scan()

			require.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestScanCommentToken(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken []Token
	}{
		{
			name:         "Valid comment with N open brackets",
			preparedText: "{{{ab}",
			expectedToken: []Token{
				NewToken(COMMENT, "{{{ab}", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Close comment twice with characters in between",
			preparedText: "{ab}ab}",
			expectedToken: []Token{
				NewToken(COMMENT, "{ab}", NULL),
				NewToken(ERROR, "", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Close comment twice",
			preparedText: "{{abab}}",
			expectedToken: []Token{
				NewToken(COMMENT, "{{abab}", NULL),
				NewToken(ERROR, "", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Comment not closed",
			preparedText: "{{abab",
			expectedToken: []Token{
				NewToken(ERROR, "", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Comment not open",
			preparedText: "abab}}",
			expectedToken: []Token{
				NewToken(ERROR, "", NULL),
				NewToken(ERROR, "", NULL),
				NewToken(EOF, "", NULL)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())

			for _, expectedToken := range tc.expectedToken {
				token := scanner.Scan()
				require.Equal(t, expectedToken, token)
			}
		})
	}
}

func TestScanLiteralConstantToken(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken Token
	}{
		{
			name:          "Simple Constant Literal",
			preparedText:  `"This is a constant literal"`,
			expectedToken: NewToken(LITERAL_CONST, `"This is a constant literal"`, LITERAL),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())
			token := scanner.Scan()

			require.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestScanGeneralCases(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken []Token
	}{
		{
			name:         "Assignment",
			preparedText: "A<-B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ATTR, "<-", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Assignment with sum",
			preparedText: "A<-B+C",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ATTR, "<-", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Assignment with subtraction",
			preparedText: "A<-B-C",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ATTR, "<-", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(ARIT_OP, "-", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Less than or greater than",
			preparedText: "A<>B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Less than or equal",
			preparedText: "A<=B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "<=", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Greater than or equal",
			preparedText: "A>=B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, ">=", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Equal",
			preparedText: "A=B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "=", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Less than",
			preparedText: "A<B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "<", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Greater than",
			preparedText: "A>B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, ">", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Operation with comparison between parentheses",
			preparedText: "(A+B<>C)",
			expectedToken: []Token{
				NewToken(OPEN_PAR, "(", NULL),
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(CLOSE_PAR, ")", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Two Operations with comparisons between parentheses and semicolon",
			preparedText: "(A+B<>C/D);",
			expectedToken: []Token{
				NewToken(OPEN_PAR, "(", NULL),
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(ARIT_OP, "/", NULL),
				NewToken(IDENTIFIER, "D", NULL),
				NewToken(CLOSE_PAR, ")", NULL),
				NewToken(SEMICOLON, ";", NULL),
				NewToken(EOF, "", NULL)},
		},
		{
			name:         "Two Operations with comparisons between parentheses and semicolon",
			preparedText: "se(A+B<>C/D);",
			expectedToken: []Token{
				NewToken("se", "se", "se"),
				NewToken(OPEN_PAR, "(", NULL),
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(ARIT_OP, "/", NULL),
				NewToken(IDENTIFIER, "D", NULL),
				NewToken(CLOSE_PAR, ")", NULL),
				NewToken(SEMICOLON, ";", NULL),
				NewToken(EOF, "", NULL)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())

			for _, expectedToken := range tc.expectedToken {
				token := scanner.Scan()
				require.Equal(t, expectedToken, token)
			}
		})
	}
}
