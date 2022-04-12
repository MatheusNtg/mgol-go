package lexer

import (
	"fmt"
	"strings"
)

type TokenClass string

// Available classes of tokens
const (
	NUM           TokenClass = "Num"
	LITERAL_CONST TokenClass = "Lit"
	IDENTIFIER    TokenClass = "id"
	COMMENT       TokenClass = "Comentário"
	REL_OP        TokenClass = "OPR"
	ARIT_OP       TokenClass = "OPM"
	EOF           TokenClass = "EOF"
	ATTR          TokenClass = "RCB"
	OPEN_PAR      TokenClass = "AB_P"
	CLOSE_PAR     TokenClass = "FC_P"
	SEMICOLON     TokenClass = "PT_V"
	ERROR         TokenClass = "ERRO"
)

type DataType string

// Available types of data
const (
	INTEGER DataType = "inteiro"
	REAL    DataType = "real"
	LITERAL DataType = "literal"
	NULL    DataType = "NULO"
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
	COMMENT_TOKEN = Token{
		class:    COMMENT,
		lexeme:   "",
		dataType: NULL,
	}
)

//Language Reserved Tokens
var LanguageReservedTokens = []Token{
	NewToken("inicio", "inicio", "inicio"),
	NewToken("varinicio", "varinicio", "varinicio"),
	NewToken("varfim", "varfim", "varfim"),
	NewToken("escreva", "escreva", "escreva"),
	NewToken("leia", "leia", "leia"),
	NewToken("se", "se", "se"),
	NewToken("entao", "entao", "entao"),
	NewToken("fimse", "fimse", "fimse"),
	NewToken("repita", "repita", "repita"),
	NewToken("fimrepita", "fimrepita", "fimrepita"),
	NewToken("fim", "fim", "fim"),
	NewToken("inteiro", "inteiro", "inteiro"),
	NewToken("literal", "literal", "literal"),
	NewToken("real", "real", "real"),
}

func NewToken(class TokenClass, lexeme string, dataType DataType) Token {
	return Token{
		class:    class,
		lexeme:   lexeme,
		dataType: dataType,
	}
}

func (t Token) GetLexem() string {
	return t.lexeme
}

func (t Token) GetClass() string {
	return strings.ToLower(string(t.class))
}

func (t Token) GetType() DataType {
	return t.dataType
}

func (t *Token) SetType(dataType DataType) {
	t.dataType = dataType
}

func (t *Token) SetClass(class TokenClass) {
	t.class = class
}

func (t Token) String() string {
	return fmt.Sprintf("Classe: %v, Lexema: %v, Tipo: %v", t.class, t.lexeme, t.dataType)
}
