package main

import (
	"github.com/opennetworktools/NetConfigParser/internal/parser"
)

func main() {
	path := "tests/configs/running-config-1.txt" // Path to config file
	parser := parser.NewParser(path)
	parser.ParseConfig()
	// fmt.Println(parser)
	// fmt.Printf("%+v\n", parser)
}

// func (p Parser) ParseConfig() error {
// 	fmt.Println("Parsing configs...")
// 	file, err := os.Open(p.FilePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if matches := reBGP.FindStringSubmatch(line); matches != nil {
// 			for i, e := range matches {
// 				fmt.Printf("%v, %v \n", i, e)
// 			}
// 		}
// 	}
// }
