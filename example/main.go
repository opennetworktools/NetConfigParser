package main

import (
	"fmt"
	"log"

	netconfigparser "github.com/opennetworktools/NetConfigParser"
)

func main() {
	osType := "IOSXE"
	path := "running-config.txt" // Path to config file
	parser := netconfigparser.GetParser(osType, path)
	if parser == nil {
		log.Fatalf("Unsupported OS type: %s", osType)
	}
	parser.ParseConfig()
 
	err := netconfigparser.WriteConfigStructToJSON(parser.GetConfigs(), "out/config.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
	}
	fmt.Println("Configs saved as JSON to configs.json")

	// pprint
	// fmt.Printf("%+v\n", parser)
}
