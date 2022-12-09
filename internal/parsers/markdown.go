package parsers

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatal("Opening the Markdown file failed.")
		panic(e)
	}
}

func ParseMdFile(filepath string) {
	mdBytes, err := os.ReadFile(filepath)
	check(err)
	fmt.Print(string(mdBytes))

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	tree := markdown.Parse(mdBytes, parser)
	log.Println("AST parsing succeed...")
	fmt.Print(tree)
}
