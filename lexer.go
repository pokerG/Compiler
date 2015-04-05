package main

import (
	"container/list"
	"fmt"
	. "github.com/pokerG/Compiler/common"
	"strings"
)

type Lexer struct {
	input        string
	input_length int
	pos          int
	line_number  int
	char_number  int
	start_pos    int
	current_char byte
	running      bool
	buffer       string
	token_stream list.List
}

type Token struct {
	token_type  int
	content     string
	line_number int
	char_number int
}

func (self *Lexer) flushBuffer() string {
	result := self.buffer
	self.buffer = ""
	return result
}

func (self *Lexer) createToken(token_type int, content string) {
	token := &Token{}
	token.token_type = token_type
	token.content = content
	token.line_number = self.line_number
	token.char_number = self.char_number
	fmt.Printf("(type %d , content %s)\n", token.token_type, token.content)
	self.token_stream.PushBack(token)
}

func (self *Lexer) peek(ahead int) byte {
	return self.input[self.pos+ahead]
}

func (self *Lexer) createLexer(input string) {
	self.input = input
	self.pos = 0
	self.input_length = len(self.input)
	self.line_number = 1
	self.char_number = 1
	self.start_pos = 0
	self.current_char = self.input[self.pos]
	self.running = true
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isLayout(c byte) bool {
	return c <= 32
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isLetterOrDigit(c byte) bool {
	return isDigit(c) || isLetter(c)
}

func isOperator(c byte) bool {
	return strings.Contains("+-*/=", string(c))
}

func isSeparator(c byte) bool {
	return strings.Contains(";,{}()[]", string(c))
}

func (self *Lexer) consumeCharacter() {
	if self.pos == self.input_length-1 {
		self.running = false
		return
	}

	if !isLayout(self.current_char) {
		self.buffer = self.buffer + string(self.input[self.pos])
	} else if self.current_char == '\n' {
		self.line_number += 1
	}
	self.pos++
	self.current_char = self.input[self.pos]
	self.char_number = self.char_number + 1
}

func (self *Lexer) recognizeNumberToken() {
	self.consumeCharacter()
	if self.current_char == '.' {
		self.consumeCharacter()
		for isDigit(self.current_char) {
			self.consumeCharacter()
		}
	}

	for isDigit(self.current_char) {
		if self.peek(1) == '.' {
			self.consumeCharacter()
			for isDigit(self.current_char) {
				self.consumeCharacter()
			}
		}
		self.consumeCharacter()
	}

	self.createToken(NUMBER, self.flushBuffer())
}

func (self *Lexer) recognizeIdentifierToken() {
	self.consumeCharacter()
	for isLetterOrDigit(self.current_char) {
		self.consumeCharacter()
	}
	for self.current_char == '_' && isLetterOrDigit(self.peek(1)) {
		self.consumeCharacter()
		for isLetterOrDigit(self.current_char) {
			self.consumeCharacter()
		}
	}

	codes := NewCodes()
	_, ok := codes[self.buffer]
	if ok {
		self.createToken(RESERVED, self.flushBuffer())
	} else {
		self.createToken(IDENTIFIER, self.flushBuffer())
	}

}

func (self *Lexer) recognizeSeparatorToken() {
	self.consumeCharacter()
	self.createToken(SEPARATOR, self.flushBuffer())
}

func (self *Lexer) recognizeStringToken() {
	self.consumeCharacter() // eat "

	for self.current_char != '"' {
		self.consumeCharacter()
	}

	self.consumeCharacter() // eat "

	self.createToken(STRING, self.flushBuffer())
}

func (self *Lexer) recognizeCharacterToken() {
	self.consumeCharacter()

	if isLetterOrDigit(self.current_char) {
		self.consumeCharacter()
	}

	self.consumeCharacter()

	self.createToken(CHARACTER, self.flushBuffer())
}

func (self *Lexer) recognizeLogicalopToken() {
	first := self.current_char
	self.consumeCharacter()
	if first == self.current_char {
		self.consumeCharacter()
	} else if self.current_char == '=' {
		self.consumeCharacter()
	}
	self.createToken(LOGICALOP, self.flushBuffer())
}

func (self *Lexer) recognizeOperatorToken() {
	self.consumeCharacter()
	self.createToken(OPERATOR, self.flushBuffer())
}

func (self *Lexer) getNextToken() {
	self.start_pos = 0
	for isLayout(self.current_char) {
		self.consumeCharacter()
		if !self.running {
			return
		}
	}
	self.start_pos = self.pos

	if isDigit(self.current_char) || self.current_char == '.' {
		self.recognizeNumberToken()
	} else if isLetterOrDigit(self.current_char) || self.current_char == '_' {
		self.recognizeIdentifierToken()
	} else if self.current_char == '"' {
		self.recognizeStringToken()
	} else if self.current_char == '\'' {
		self.recognizeCharacterToken()
	} else if strings.Contains("!=|&<>", string(self.current_char)) {
		self.recognizeLogicalopToken()
	} else if isOperator(self.current_char) {
		self.recognizeOperatorToken()
	} else if isSeparator(self.current_char) {
		self.recognizeSeparatorToken()
	} else {
		fmt.Printf("unknown token type %d, aka %c at %d\n", self.current_char, self.current_char, self.line_number)
		self.running = false
		return
	}
}

func (self *Lexer) startLexing() {
	for self.running {
		self.getNextToken()
	}

	self.createToken(END_OF_FILE, "<EOF>")
}
