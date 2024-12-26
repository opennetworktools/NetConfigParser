package parser

import "fmt"

type EOSParser struct {
}

func (p *EOSParser) PrintParserType() {
	fmt.Println("Arista EOS")
}
