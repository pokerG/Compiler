package main

import (
	"io/ioutil"
)


var signaltable *SignalTable
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

	signaltable = NewSignalTable()
	lexer := &Lexer{}
	lexer.createLexer(file_contents)
	lexer.startLexing()
	parser := &Parser{}
	parser.createParser(lexer.token_stream)
	parser.startParsing()

}
