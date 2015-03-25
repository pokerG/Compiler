package main

import (
	"fmt"
	. "github.com/pokerG/Compiler/common"
	// "io"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	IsnotExist   = errors.New("Not have this stynax")
	InitalErr    = errors.New("Standing initial error")
	Unrecognized = errors.New("Variable have unrecognized character")
)

func main() {
	if len(os.Args) > 1 {
		parse(os.Args[1])
	} else {
		fmt.Println("Please input the file path")
	}
}

func parse(filepath string) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lines := strings.Split(string(buf), "\n")
	// x, y := 1, 1 // column and line
	// fmt.Println(strings.Split(string(buf), "\n"))
	for _, s := range lines {
		word := ""
		s = strings.TrimSpace(s)
		for forward := 0; forward <= len(s); forward++ {
			if forward == len(s) {
				if len(word) != 0 {
					Print(word)
				}

				break
			}
			if s[forward] == ' ' {
				if len(word) > 0 {
					Print(word)
				}
				word = ""
			} else if s[forward] == '/' && s[forward+1] == '/' {
				break
			} else if s[forward] == '+' || s[forward] == '-' || s[forward] == '*' || s[forward] == '/' || s[forward] == '(' ||
				s[forward] == ')' || s[forward] == '[' || s[forward] == ']' || s[forward] == '{' || s[forward] == '}' ||
				s[forward] == ',' || s[forward] == '\'' {
				if len(word) > 0 {
					Print(word)
				}
				word = string(byte(s[forward]))
				Print(word)
				word = ""
			} else if s[forward] == '<' || s[forward] == '>' || s[forward] == '=' || s[forward] == '!' || s[forward] == ':' {
				if len(word) > 0 {
					Print(word)
				}
				word = string(byte(s[forward]))
				if forward != len(s)-1 && s[forward+1] == '=' {
					word += string(byte(s[forward+1]))
					forward += 1
				}
				Print(word)
				word = ""
			} else if s[forward] == '&' && s[forward+1] == '&' {
				if len(word) > 0 {
					Print(word)
				}
				word = string(byte(s[forward])) + string(byte(s[forward+1]))
				forward += 1
				Print(word)
				word = ""
			} else if s[forward] == '|' && s[forward+1] == '|' {
				if len(word) > 0 {
					Print(word)
				}
				word = string(byte(s[forward])) + string(byte(s[forward+1]))
				Print(word)
				word = ""
			} else if (s[forward] >= '0' && s[forward] <= '9') || (s[forward] >= 'a' && s[forward] <= 'z') ||
				(s[forward] >= 'A' && s[forward] <= 'Z') || s[forward] == '.' || s[forward] == '_' {
				word += string(byte(s[forward]))
			}
		}
	}
}

func Print(word string) error {
	// fmt.Println("#" + word)
	codes := NewCodes()
	_, ok := codes[word]
	if ok {
		fmt.Println(word + "," + word)
		return nil
	}
	err := isNumber(word)
	if err == nil {
		fmt.Println("Constant," + word)
		return nil
	}
	err = isVariable(word)
	if err == nil {
		fmt.Println("Variable," + word)
		return nil
	}
	return IsnotExist
}

func isNumber(word string) error {
	_, err := strconv.Atoi(word)
	return err
}

func isVariable(word string) error {

	for k, v := range word {
		if k == 1 {
			if v != '_' && (v < 'A' || v > 'Z') &&
				(v < 'a' || v > 'z') {
				return InitalErr
			}
		}
		if v != '_' && (v < 'A' || v > 'Z') &&
			(v < 'a' || v > 'z') && (v < '0' || v > '9') {
			return Unrecognized
		}
	}
	return nil
}
