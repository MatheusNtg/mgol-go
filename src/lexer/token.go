package lexer

type TokenClass int

// Available classes of tokens
const (
	NUM TokenClass = iota
	LITERAL_CONST
	IDENTIFIER
	COMMENT
	REL_OP
	ARIT_OP
	EOF
	ATTR
	OPEN_PAR
	CLOSE_PAR
	SEMICOLON
	ERROR
)

type DataType int

// Available types of data
const (
	INTEGER DataType = iota
	REAL
	LITERAL
	NULL
)

type Token struct {
	class    TokenClass
	lexeme   string
	dataType DataType
}

// Constant Tokens
var (
	EOF_TOKEN = Token{
		class:    EOF,
		lexeme:   "",
		dataType: NULL,
	}
	ATTR_TOKEN = Token{
		class:    ATTR,
		lexeme:   "<-",
		dataType: NULL,
	}
	OPEN_PAR_TOKEN = Token{
		class:    OPEN_PAR,
		lexeme:   "(",
		dataType: NULL,
	}
	CLOSE_PAR_TOKEN = Token{
		class:    CLOSE_PAR,
		lexeme:   ")",
		dataType: NULL,
	}
	SEMICOLON_TOKEN = Token{
		class:    SEMICOLON,
		lexeme:   ";",
		dataType: NULL,
	}
	ERROR_TOKEN = Token{
		class:    ERROR,
		lexeme:   "",
		dataType: NULL,
	}
)

func NewToken(class TokenClass, lexeme string, dataType DataType) Token {
	return Token{
		class:    class,
		lexeme:   lexeme,
		dataType: dataType,
	}
}
