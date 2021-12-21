package lexer

type TokenClass int

// Available classes of tokens
const (
	NUM TokenClass = iota
	LITERAL_CONST
	IDENTIFIER
	COMMENT
	EOF
	REL_OP
	ATTR
	ARIT_OP
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

func NewToken(class TokenClass, lexeme string, dataType DataType) *Token {
	return &Token{
		class:    class,
		lexeme:   lexeme,
		dataType: dataType,
	}
}
