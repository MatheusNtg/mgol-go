package symboltable

import (
	"mgol-go/src/lexer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError error
		keys          []string
		values        []lexer.Token
		errorIndex    int
	}{
		{
			name:          "Insert without conflict",
			expectedError: nil,
			keys:          []string{"k1", "k2", "k3"},
			values: []lexer.Token{
				lexer.NewToken(lexer.COMMENT, "comment", lexer.NULL),
				lexer.NewToken(lexer.IDENTIFIER, "identi", lexer.NULL),
				lexer.NewToken(lexer.LITERAL_CONST, `"an test"`, lexer.LITERAL),
			},
		},
		{
			name:          "Insert with conflict",
			expectedError: ErrorAlreadyOnTable,
			keys:          []string{"k1", "k1", "k3"},
			values: []lexer.Token{
				lexer.NewToken(lexer.COMMENT, "comment", lexer.NULL),
				lexer.NewToken(lexer.IDENTIFIER, "identi", lexer.NULL),
				lexer.NewToken(lexer.LITERAL_CONST, `"an test"`, lexer.LITERAL),
			},
			errorIndex: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for index, key := range tc.keys {
				token, err := Insert(key, tc.values[index])
				if tc.expectedError != nil && index == tc.errorIndex {
					require.Error(t, err)
					continue
				}
				require.NoError(t, err)
				require.Equal(t, tc.values[index], token)
			}
			CleanupTable()
		})
	}
}

func TestGetToken(t *testing.T) {
	testCases := []struct {
		name            string
		expectedError   error
		prepareFunction func()
		key             string
		expectedToken   lexer.Token
	}{
		{
			name:          "Get existing token",
			expectedError: nil,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k2", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k2",
			expectedToken: lexer.CLOSE_PAR_TOKEN,
		},
		{
			name:          "Get non-existing token",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k2", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k7",
			expectedToken: lexer.Token{},
		},
		{
			name:          "Get existing token on confliting table",
			expectedError: nil,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k1", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k1",
			expectedToken: lexer.ATTR_TOKEN,
		},
		{
			name:          "Get non-existing token on confliting table",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k1", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k7",
			expectedToken: lexer.Token{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareFunction()
			token, err := GetToken(tc.key)
			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedToken, token)
			CleanupTable()
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name            string
		expectedError   error
		prepareFunction func()
		key             string
		newToken        lexer.Token
		expectedToken   lexer.Token
	}{
		{
			name:          "Successfully update",
			expectedError: nil,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k2", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k2",
			newToken:      lexer.EOF_TOKEN,
			expectedToken: lexer.EOF_TOKEN,
		},
		{
			name:          "Update an non-existing token",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k2", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k4",
			newToken:      lexer.Token{},
			expectedToken: lexer.Token{},
		},
		{
			name:          "Successfully update on conflict tables",
			expectedError: nil,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k2", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k2",
			newToken:      lexer.EOF_TOKEN,
			expectedToken: lexer.EOF_TOKEN,
		},
		{
			name:          "Update an non-existing token on conflict tables",
			expectedError: ErrorSymbolNotFound,
			prepareFunction: func() {
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k1", lexer.ATTR_TOKEN)
				Insert("k2", lexer.CLOSE_PAR_TOKEN)
				Insert("k3", lexer.OPEN_PAR_TOKEN)
			},
			key:           "k4",
			newToken:      lexer.Token{},
			expectedToken: lexer.Token{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareFunction()
			err := Update(tc.key, tc.newToken)
			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			table := GetTable()
			require.Equal(t, tc.expectedToken, table[tc.key])
			CleanupTable()
		})
	}
}
