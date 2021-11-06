package main

import (
	"fmt"

	"github.com/kingzbauer/scraperlang/cmdutil"
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
	fmt.Println(ast)
}
