package main

import (
	"log"
	"strconv"
	"unicode/utf8"
)

// reserved words
var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *scanner {
	s := scanner{source: source, start: 0, current: 0, line: 0}
	return &s
}

func (code *scanner) scanTokens() {
	for {
		if code.isAtEnd() {
			break
		}
		code.scanToken()
	}
}

func (code *scanner) isAtEnd() bool {
	return code.current >= len(code.source)
}

func (code *scanner) advance() rune {
	if code.isAtEnd() {
		return 0
	}
	r, size := utf8.DecodeRuneInString(code.source[code.current:])
	code.current += size
	return r
}

func (code *scanner) addToken(token_type TokenType) {
	code.addTokenValue(token_type, nil)
}

func (code *scanner) addTokenValue(token_type TokenType, literal any) {
	text := code.source[code.start:code.current]
	code.tokens = append(code.tokens, Token{token_type: token_type, lexeme: text, literal: literal, line: code.line})
	code.start = code.current
}

func (code *scanner) match(expected byte) bool {
	if code.isAtEnd() {
		return false
	}
	if code.source[code.current] != expected {
		return false
	}
	code.current++
	return true
}

func (code *scanner) scanToken() {
	c := code.advance()
	switch c {
	case '(':
		code.addToken(LEFT_PAREN)
	case ')':
		code.addToken(RIGHT_PAREN)
	case '{':
		code.addToken(LEFT_BRACE)
	case '}':
		code.addToken(RIGHT_BRACE)
	case ',':
		code.addToken(COMMA)
	case '.':
		code.addToken(DOT)
	case '-':
		code.addToken(MINUS)
	case '+':
		code.addToken(PLUS)
	case ';':
		code.addToken(SEMICOLON)
	case '*':
		code.addToken(STAR)
	case '!':
		if code.match('=') {
			code.addToken(BANG_EQUAL)
		} else {
			code.addToken(BANG)
		}
	case '=':
		if code.match('=') {
			code.addToken(EQUAL_EQUAL)
		} else {
			code.addToken(EQUAL)
		}
	case '<':
		if code.match('=') {
			code.addToken(LESS_EQUAL)
		} else {
			code.addToken(LESS)
		}
	case '>':
		if code.match('=') {
			code.addToken(GREATER_EQUAL)
		} else {
			code.addToken(GREATER)
		}
	case '/':
		if code.match('/') {
			for (code.peek() != '\n') && (!code.isAtEnd()) {
				code.advance()
			}
			//code.start = code.current
			code.nop()
		} else {
			code.addToken(SLASH)
		}
	case '\n':
		code.line++
		code.nop()
	case '\r', '\t', ' ':
		// ignore white spaces
		code.nop()
	case '"':
		code.read_string()
	default:
		if code.isDigit(c) {
			code.number()
		} else if code.isAlpha(c) {
			code.identifier()
		} else {
			errorReport(code.line, "Unexpected character. |"+strconv.Itoa(int(c))+"|")
		}
	}
}

func (code *scanner) nop() {
	code.start = code.current
}

func (code *scanner) peek() rune {
	if code.isAtEnd() {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(code.source[code.current:])
	return r
}

func (code *scanner) peekNext() rune {
	if (code.current + 1) >= len(code.source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(code.source[code.current+1:])
	return r
}

func (code *scanner) read_string() {
	for (code.peek() != '"') && (!code.isAtEnd()) {
		if code.peek() == '\n' {
			code.line++
		}
		code.advance()
	}
	if code.isAtEnd() {
		errorReport(code.line, "Unterminated string")
	}
	code.advance()
	text := code.source[code.start+1 : code.current-1]
	code.addTokenValue(STRING, text)
}

func (code *scanner) isDigit(c rune) bool {
	return (c >= '0') && (c <= '9')
}

func (code *scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}
func (code *scanner) isAlphanumeric(c rune) bool {
	return code.isAlpha(c) || code.isDigit(c)
}

func (code *scanner) number() {
	for code.isDigit(code.peek()) {
		code.advance()
	}

	if code.peek() == '.' && code.isDigit(code.peekNext()) {
		code.advance()
		for code.isDigit(code.peek()) {
			code.advance()
		}
	}
	valorStr := code.source[code.start:code.current]
	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		log.Println("Unable to parse " + valorStr + " into float")
	}
	code.addTokenValue(NUMBER, valor)
}

func (code *scanner) identifier() {
	for code.isAlphanumeric(code.peek()) {
		code.advance()
	}

	text := code.source[code.start:code.current]
	token_type, ok := keywords[text]
	if !ok {
		token_type = IDENTIFIER
	}
	code.addToken(token_type)
}
