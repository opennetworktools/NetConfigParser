package main

import (
	"fmt"

	"github.com/opennetworktools/NetConfigParser/internal/parser"
	"github.com/opennetworktools/NetConfigParser/internal/utils"
)

func main() {
	path := "tests/configs/running-config-1.txt" // Path to config file
	parser := parser.NewParser(path)
	parser.ParseConfig()

	// preety-print
	// fmt.Printf("%+v\n", parser)

	err := utils.WriteParserStructToJSON(parser, "out/parser.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
	}
	fmt.Println("Parser saved as JSON to parser.json")
}
