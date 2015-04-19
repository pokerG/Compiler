package main

import (
	"io/ioutil"
)

func readFile(fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func main() {
	// read the file contents
	file_contents := readFile("test")

	var eof byte = 0
	file_contents += string(eof)

	// load up our lexer
	lexer := &Lexer{}
	lexer.createLexer(file_contents)
	lexer.startLexing()
	// parser := &Parser{}
	// parser.createParser(lexer.token_stream)

}
