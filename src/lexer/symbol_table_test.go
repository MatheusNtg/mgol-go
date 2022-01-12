package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	testCases := []struct {
		name           string
		keys           []string
		values         []Token
		expectedResult []Token
	}{
		{
			name: "Insert without conflict",
			keys: []string{"k1", "k2", "k3"},
			values: []Token{
				NewToken(COMMENT, "comment", NULL),
				NewToken(IDENTIFIER, "identi", NULL),
				NewToken(LITERAL_CONST, `"an test"`, LITERAL),
			},
			expectedResult: []Token{
				NewToken(COMMENT, "comment", NULL),
				NewToken(IDENTIFIER, "identi", NULL),
				NewToken(LITERAL_CONST, `"an test"`, LITERAL),
			},
		},
		{
			name: "Insert with conflict",
			keys: []string{"k1", "k1", "k3"},
			values: []Token{
				NewToken(COMMENT, "comment", NULL),
				NewToken(IDENTIFIER, "identi", NULL),
				NewToken(LITERAL_CONST, `"an test"`, LITERAL),
			},
			expectedResult: []Token{
				NewToken(COMMENT, "comment", NULL),
				NewToken(COMMENT, "comment", NULL),
				NewToken(LITERAL_CONST, `"an test"`, LITERAL),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for index, key := range tc.keys {
				token := InsertSymbolTable(key, tc.values[index])
				require.Equal(t, tc.expectedResult[index], token)
			}
			CleanupSymbolTable()
		})
	}
}

func TestGetToken(t *testing.T) {
	testCases := []struct {
		name            string
		expectedError   error
		prepareFunction func()
		key             string
		expectedToken   Token
	}{
		{
			name:          "Get existing token",
			expectedError: nil,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k2", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k2",
			expectedToken: CLOSE_PAR_TOKEN,
		},
		{
			name:          "Get non-existing token",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k2", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k7",
			expectedToken: Token{},
		},
		{
			name:          "Get existing token on confliting table",
			expectedError: nil,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k1", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k1",
			expectedToken: ATTR_TOKEN,
		},
		{
			name:          "Get non-existing token on confliting table",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k1", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k7",
			expectedToken: Token{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareFunction()
			token, err := GetTokenFromSymbolTable(tc.key)
			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedToken, token)
			CleanupSymbolTable()
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name            string
		expectedError   error
		prepareFunction func()
		key             string
		newToken        Token
		expectedToken   Token
	}{
		{
			name:          "Successfully update",
			expectedError: nil,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k2", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k2",
			newToken:      EOF_TOKEN,
			expectedToken: EOF_TOKEN,
		},
		{
			name:          "Update an non-existing token",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k2", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k4",
			newToken:      Token{},
			expectedToken: Token{},
		},
		{
			name:          "Successfully update on conflict tables",
			expectedError: nil,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k2", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k2",
			newToken:      EOF_TOKEN,
			expectedToken: EOF_TOKEN,
		},
		{
			name:          "Update an non-existing token on conflict tables",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k1", ATTR_TOKEN)
				InsertSymbolTable("k2", CLOSE_PAR_TOKEN)
				InsertSymbolTable("k3", OPEN_PAR_TOKEN)
			},
			key:           "k4",
			newToken:      Token{},
			expectedToken: Token{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareFunction()
			err := UpdateSymbolTable(tc.key, tc.newToken)
			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			table := GetSymbolTable()
			require.Equal(t, tc.expectedToken, table[tc.key])
			CleanupSymbolTable()
		})
	}
}
