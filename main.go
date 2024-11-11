package main

import (
	"fmt"
	"log"

	"github.com/opennetworktools/NetConfigParser/internal/parser"
	"github.com/opennetworktools/NetConfigParser/internal/utils"
)

func main() {
	osType := "IOSXE"
	path := "tests/configs/running-config-1.txt" // Path to config file
	parser := parser.GetParser(osType, path)
	if parser == nil {
		log.Fatalf("Unsupported OS type: %s", osType)
	}
	parser.ParseConfig()

	// pprint
	// fmt.Printf("%+v\n", parser)
 
	err := utils.WriteParserStructToJSON(parser.GetConfigs(), "out/config.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
	}
	fmt.Println("Configs saved as JSON to configs.json")
 }
 
 
