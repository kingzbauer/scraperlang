package main

import (
	"fmt"
	"os"

	"github.com/kingzbauer/scraperlang/cmdutil"
	"github.com/kingzbauer/scraperlang/interpreter"
	"github.com/kingzbauer/scraperlang/parser"
	"github.com/kingzbauer/scraperlang/token"
)

func main() {
	src, err := cmdutil.ReadFileArg(true)
	cmdutil.ExitOnError(err)

	scanner := token.NewScanner(src)
	tokens, err := scanner.ScanTokens()
	cmdutil.ExitOnError(err)

	p := parser.New(tokens)
	ast, err := p.Parse()
	cmdutil.ExitOnError(err)
	if p.HasErrs() {
		for _, err := range p.Err() {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	i, err := interpreter.New(ast)
	cmdutil.ExitOnError(err)
	cmdutil.ExitOnError(i.Exec())
}
