package main

import (
	"fmt"
	"os"

	"github.com/kingzbauer/scraperlang/token"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Missing filename")
		os.Exit(1)
	}
	filename := os.Args[1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	scanner := token.NewScanner(content)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v", tokens)
}
