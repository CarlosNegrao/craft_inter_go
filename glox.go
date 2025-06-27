package main

import (
	"fmt"
	"strings"

	//"unicode/utf8"
	"bufio"
	"io"
	"log"
	"os"
)

// variáveis do pacote
var hadError bool = false

// para facilitar tratamento de erros
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 1 {
		fmt.Println("Usage: jlox [script]")
	} else if len(argsWithoutProg) == 1 {
		runFile(argsWithoutProg[0])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	dat, err := os.ReadFile(path)
	check(err)
	run(string(dat))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	leitor := bufio.NewReader(os.Stdin)
	for {
		texto, err := leitor.ReadString('\n')
		if err == io.EOF {
			break
		}
		run(texto)
		hadError = false
	}
}

func run(source string) {
	scanner := NewScanner(strings.Trim(source, "\n\r"))
	scanner.scanTokens()

	for i := range len(scanner.tokens) {
		fmt.Println(scanner.tokens[i])
	}
}

func errorReport(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Println(line, where, message)
	hadError = true
}
