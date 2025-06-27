package main

import (
	"fmt"
)

type Token struct {
	token_type TokenType
	lexeme     string
	literal    any
	line       int
}

// Implementa fmt.Stringer
func (t Token) String() string {
	return fmt.Sprintf("[type=%d lexeme=%q literal=%v line=%d]", t.token_type, t.lexeme, t.literal, t.line)
}
