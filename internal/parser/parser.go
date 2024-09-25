package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/opennetworktools/NetConfigParser/internal/net"
)

type Parser struct {
	FilePath string
	Configs  Configs
}

type Configs struct {
	BGP        net.BGP
	Interfaces net.Interfaces
}

func NewParser(filePath string) *Parser {
	return &Parser{
		FilePath: filePath,
	}
}

func (p *Parser) ParseConfig() error {
	fmt.Println("Parsing configs...")

	content, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	configString := string(content)

	configs := Configs{}

	// BGP block
	reBGP := regexp.MustCompile(`(?s)router bgp.*?!\n`)
	bgpBlock := reBGP.FindString(configString)
	bgp_object := net.ParseBGPBlock(bgpBlock)
	configs.BGP = bgp_object
	p.Configs = configs

	// All Interfaces block
	// reInterface := regexp.MustCompile(`(?s)(interface .+?)(?=^!|$)`)
	// reInterface := regexp.MustCompile(`(?m)^interface \S+`)
	reInterface := regexp.MustCompile(`^interface`)
	matches := reInterface.FindAllString(configString, -1)
	fmt.Println(matches)

	// var blocks []string
	// for i, match := range matches {
	// 	start := match[0]
	// 	var end int
	// 	if i+1 < len(matches) {
	// 		end = match[1][0]
	// 	} else {
	// 		end = len(configString)
	// 	}
	// 	blocks = append(blocks, strings.TrimSpace(configString[start:end]))
	// }

	for i, match := range matches {
		fmt.Println(i)
		fmt.Println(match)
		// fmt.Println("!")
	}

	// fmt.Println(blocks)

	return nil
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