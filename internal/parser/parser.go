package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

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

	// InterfacesBlock
	regexInterface := `^\s*interface .+$`
	reInterface := regexp.MustCompile(regexInterface)
	interfacesBlock, err := extractBlock(configString, reInterface)
	if err != nil {
		fmt.Println("Error extracting interfaces block:", err)
		return err
	}
	interfacesObj := net.ParseInterfacesBlock(interfacesBlock)
	// fmt.Println(interfacesObj)
	configs.Interfaces = interfacesObj

	// Set configs to parser object
	p.Configs = configs

	return nil
}

func extractBlock(s string, re *regexp.Regexp) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	var block, blocks []string
	for scanner.Scan() {
		line := scanner.Text()
		if block == nil && re.MatchString(line) {
			block = append(block, line)
		} else if block != nil && strings.HasPrefix(line, " ") {
			block = append(block, line)
		} else if block != nil && re.MatchString(line) {
			blocks = append(blocks, strings.Join(block, "\n"))
			block = nil
			block = append(block, line)
		} else if block != nil && !strings.HasPrefix(line, " ") {
			blocks = append(blocks, strings.Join(block, "\n"))
			block = nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return blocks, nil
}

func PrettyPrint(arr []string) {
	for i, item := range arr {
		fmt.Println(i)
		fmt.Println(item)
	}
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