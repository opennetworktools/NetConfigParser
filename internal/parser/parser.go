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
	BGP           net.BGP
	Interfaces    net.Interfaces
	RouteMaps     []net.RouteMap
	IPAccessLists []net.AccessList
	IPPrefixLists []net.PrefixList
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
	interfacesBlock, err := extractBlock(configString, reInterface, "OTHER")
	if err != nil {
		fmt.Println("Error extracting interfaces block:", err)
		return err
	}
	interfacesObj := net.ParseInterfacesBlock(interfacesBlock)
	configs.Interfaces = interfacesObj

	// Route-map
	regexRouteMap := `^\s*route-map .+$`
	reRouteMap := regexp.MustCompile(regexRouteMap)
	routeMapsBlock, err := extractBlock(configString, reRouteMap, "OTHER")
	if err != nil {
		fmt.Println("Error extracting route-map block:", err)
		return err
	}
	routeMapsObj := net.ParseRouteMapBlock(routeMapsBlock)
	configs.RouteMaps = routeMapsObj

	// IP Access-list
	regexIPAccessList := `^\s*ip access-list .+$`
	reIPAccessList := regexp.MustCompile(regexIPAccessList)
	IPAccessListsBlock, err := extractBlock(configString, reIPAccessList, "OTHER")
	if err != nil {
		fmt.Println("Error extracting ip access-list block:", err)
		return err
	}
	IPAccessListObj := net.ParseIPAccessListBlock(IPAccessListsBlock)
	configs.IPAccessLists = IPAccessListObj

	// IP Prefix-list
	regexIPPrefixList := `^\s*ip prefix-list .+$`
	reIPPrefixList := regexp.MustCompile(regexIPPrefixList)
	IPPrefixListsBlock, err := extractBlock(configString, reIPPrefixList, "IP_PREFIX_LIST")
	if err != nil {
		fmt.Println("Error extracting ip access-list block:", err)
		return err
	}
	IPPrefixListObj := net.ParseIPPrefixListBlock(IPPrefixListsBlock)
	configs.IPPrefixLists = IPPrefixListObj

	// Set configs to parser object
	p.Configs = configs

	return nil
}

func extractBlock(s string, re *regexp.Regexp, feature string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	var block, blocks []string
	count := 1
	currentLine := ""
	for scanner.Scan() {
		count++;
		prevLine := currentLine
		currentLine = scanner.Text()
		// line := scanner.Text()
		if block == nil && re.MatchString(currentLine) {
			// first instance match of the line
			block = append(block, currentLine)
		} else if block != nil && strings.HasPrefix(currentLine, " ") {
			// if the line has prefix of space character
			block = append(block, currentLine)
		} else if block != nil && re.MatchString(currentLine) {
			// not the first instance match of the line
			if prevLine != "!" && feature == "IP_PREFIX_LIST" {
				block = append(block, currentLine)
				continue
			}
			blocks = append(blocks, strings.Join(block, "\n"))
			block = nil
			block = append(block, currentLine)
		} else if block != nil && !strings.HasPrefix(currentLine, " ") {
			// if the line doesn't has prefix of space character
			blocks = append(blocks, strings.Join(block, "\n"))
			block = nil
		}
	}
	// fmt.Printf("extractBlock function count: %v \n", count)
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