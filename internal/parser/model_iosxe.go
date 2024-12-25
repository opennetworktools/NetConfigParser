package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/opennetworktools/NetConfigParser/internal/net"
)

type IOSXEParser struct {
   FilePath string
   Features []string
   Configs  Configs
}

func (p *IOSXEParser) GetConfigs() *Configs {
   return &p.Configs
}

func (p *IOSXEParser) GetFeatures() *[]string {
    return &p.Features
}

func (p *IOSXEParser) PrintParserType() {
   fmt.Println("Cisco IOS-XE")
}

func (p *IOSXEParser) ParseConfig() error {
   fmt.Println("Parsing configs...")

   content, err := ioutil.ReadFile(p.FilePath)
   if err != nil {
       fmt.Println("Error reading file:", err)
       return err
   }

   configString := string(content)

   configs := Configs{}

   // Metadata
   metadataObject := p.ParseMetadata(configString)
   configs.Metadata = metadataObject

   // BGP block
   reBGP := regexp.MustCompile(`(?s)router bgp.*?!\n`)
   bgpBlock := reBGP.FindString(configString)
   bgpObject := p.ParseBGPBlock(bgpBlock)
   configs.BGP = bgpObject
   if bgpObject.ASN != "" {
       p.Features = append(p.Features, "BGP")
   }

   // InterfacesBlock
   regexInterface := `^\s*interface .+$`
   reInterface := regexp.MustCompile(regexInterface)
   interfacesBlock, err := extractBlock(configString, reInterface, "OTHER")
   if err != nil {
       fmt.Println("Error extracting interfaces block:", err)
       return err
   }
   interfacesObj := p.ParseInterfacesBlock(interfacesBlock)
   configs.Interfaces = interfacesObj

   // Route-map
   regexRouteMap := `^\s*route-map .+$`
   reRouteMap := regexp.MustCompile(regexRouteMap)
   routeMapsBlock, err := extractBlock(configString, reRouteMap, "OTHER")
   if err != nil {
       fmt.Println("Error extracting route-map block:", err)
       return err
   }
   routeMapsObj := p.ParseRouteMapBlock(routeMapsBlock)
   configs.RouteMaps = routeMapsObj
   if len(routeMapsObj) > 0 {
    p.Features = append(p.Features, "Route-Map")
   }

   // IP Access-list
   regexIPAccessList := `^\s*ip access-list .+$`
   reIPAccessList := regexp.MustCompile(regexIPAccessList)
   IPAccessListsBlock, err := extractBlock(configString, reIPAccessList, "OTHER")
   if err != nil {
       fmt.Println("Error extracting ip access-list block:", err)
       return err
   }
   IPAccessListObj := p.ParseIPAccessListBlock(IPAccessListsBlock)
   configs.IPAccessLists = IPAccessListObj
   if len(IPAccessListObj) > 0 {
    p.Features = append(p.Features, "IP Access List")
   }

   // IP Prefix-list
	regexIPPrefixList := `^\s*ip prefix-list .+$`
	reIPPrefixList := regexp.MustCompile(regexIPPrefixList)
	IPPrefixListsBlock, err := extractBlock(configString, reIPPrefixList, "IP_PREFIX_LIST")
	if err != nil {
		fmt.Println("Error extracting ip access-list block:", err)
		return err
	}
	IPPrefixListObj := p.ParseIPPrefixListBlock(IPPrefixListsBlock)
	configs.IPPrefixLists = IPPrefixListObj
    if len(IPPrefixListObj) > 0 {
        p.Features = append(p.Features, "Prefix List")
    }

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

func (p *IOSXEParser) ParseBGPBlock(block string) net.BGP {
   fmt.Println("Parsing BGP...")
   scanner := bufio.NewScanner(strings.NewReader(block))

   // Regular expressions to match different BGP details
   reBGP := regexp.MustCompile(`^router bgp (\d+)`)
   reBGPDetail := regexp.MustCompile(`^ bgp router-id (\S+)`)

   // TODO: use IPv4 or v6 address regex pattern instead "(\S+)""
   reBGPNeighbor := regexp.MustCompile(`^ neighbor (\S+) remote-as (\d+)`)
   // reBGPNeighborDetail := regexp.MustCompile(`^ neighbor (\S+) (\S+) (.+)`)
   reBGPNeighborDescription := regexp.MustCompile(`^ neighbor (\S+) description (\S+)`)
   reBGPNeighborTimers := regexp.MustCompile(`^ neighbor (\S+) timers (\d+) (\d+)`)

   bgpObject := net.BGP{
       ASN: "0",
   }

   // Build neighbors array
   for scanner.Scan() {
       line := scanner.Text()
       if matches := reBGP.FindStringSubmatch(line); matches != nil {
           bgpObject.ASN = matches[1]
           continue
       }
       if matches := reBGPDetail.FindStringSubmatch(line); matches != nil {
           bgpObject.RouterID = matches[1]
           continue
       }
       if matches := reBGPNeighbor.FindStringSubmatch(line); matches != nil {
           bgpNeighborObject := net.BGPNeighbor{
               IP:  matches[1],
               ASN: matches[2],
           }
           bgpNeighborObject.SessionAttributes = make(map[string]string)
           bgpObject.Neighbors = append(bgpObject.Neighbors, bgpNeighborObject)
           continue
       }
       if matches := reBGPNeighborDescription.FindStringSubmatch(line); matches != nil {
           for i, neighbor := range bgpObject.Neighbors {
               if neighbor.IP == matches[1] {
                   bgpObject.Neighbors[i].SessionAttributes["Description"] = matches[2]
               }
           }
           continue
       }
       if matches := reBGPNeighborTimers.FindStringSubmatch(line); matches != nil {
           for i, neighbor := range bgpObject.Neighbors {
               if neighbor.IP == matches[1] {
                   bgpObject.Neighbors[i].SessionAttributes["HelloTimer"] = matches[2]
                   bgpObject.Neighbors[i].SessionAttributes["DeadTimer"] = matches[3]
               }
           }
           continue
       }
   }
   return bgpObject
}

func (p *IOSXEParser) ParseInterfacesBlock(blocks []string) net.Interfaces {
   fmt.Println("Parsing Interfaces Block...")
   interfaces := net.Interfaces{
       GigabitEthernet:    []net.Interface{},
       TenGigabitEthernet: []net.Interface{},
       PortChannel:        []net.Interface{},
       Loopback:           []net.Interface{},
       Tunnel:             []net.Interface{},
   }

   // Interface types REGEX
   // regexEthernet := `^\s*interface .+Ethernet\d+`
   regexEthernet := `^\s*interface Ethernet\d+`
   regexVlan := `^\s*interface Vlan\d+`
   regexFastEthernet := `^\s*interface FastEthernet\d+`
   regexGigabitEthernet := `^\s*interface GigabitEthernet\d+`
   regexTenGigabitEthernet := `^\s*interface TenGigabitEthernet\d+`
   regexLoopback := `^\s*interface Loopback\d+`
   regexPortChannel := `^\s*interface Port-channel\d+`
   regexTunnel := `^\s*interface Tunnel\d+`

   reEthernet := regexp.MustCompile(regexEthernet)
   reVlan := regexp.MustCompile(regexVlan)
   reFastEthernet := regexp.MustCompile(regexFastEthernet)
   reGigabitEthernet := regexp.MustCompile(regexGigabitEthernet)
   reTenGigabitEthernet := regexp.MustCompile(regexTenGigabitEthernet)
   reLoopback := regexp.MustCompile(regexLoopback)
   rePortChannel := regexp.MustCompile(regexPortChannel)
   reTunnel := regexp.MustCompile(regexTunnel)

   var interfaceObj net.Interface
   var interfaceSubstrings []string

   for _, block := range blocks {
       if reEthernet.MatchString(block) {
           fmt.Println("Interface Ethernet Matched!")
       } else if reVlan.MatchString(block) {
           fmt.Println("Interface Vlan Matched!")
       } else if reFastEthernet.MatchString(block) {
           fmt.Println("Interface FastEthernet Matched!")
       } else if reGigabitEthernet.MatchString(block) {
           fmt.Println("Interface GigabitEthernet Matched!")
           interfaceSubstrings = reGigabitEthernet.FindStringSubmatch(block)
           interfaceObj = p.ParseInterfaceString(block, interfaceSubstrings[0])
           interfaces.GigabitEthernet = append(interfaces.GigabitEthernet, interfaceObj)
           // fmt.Println(interfaceObj)
       } else if reTenGigabitEthernet.MatchString(block) {
           fmt.Println("Interface TenGigabitEthernet Matched!")
           interfaceSubstrings = reTenGigabitEthernet.FindStringSubmatch(block)
           interfaceObj = p.ParseInterfaceString(block, interfaceSubstrings[0])
           interfaces.TenGigabitEthernet = append(interfaces.TenGigabitEthernet, interfaceObj)
       } else if reLoopback.MatchString(block) {
           fmt.Println("Interface Loopback Matched!")
           interfaceSubstrings = reLoopback.FindStringSubmatch(block)
           interfaceObj = p.ParseInterfaceString(block, interfaceSubstrings[0])
           interfaces.Loopback = append(interfaces.Loopback, interfaceObj)
       } else if rePortChannel.MatchString(block) {
           fmt.Println("Interface PortChannel Matched!")
           interfaceSubstrings = rePortChannel.FindStringSubmatch(block)
           interfaceObj = p.ParseInterfaceString(block, interfaceSubstrings[0])
           interfaces.PortChannel = append(interfaces.PortChannel, interfaceObj)
       } else if reTunnel.MatchString(block) {
           fmt.Println("Interface Tunnel Matched!")
           interfaceSubstrings = reTunnel.FindStringSubmatch(block)
           interfaceObj = p.ParseInterfaceString(block, interfaceSubstrings[0])
           interfaces.Tunnel = append(interfaces.Tunnel, interfaceObj)
       }
       interfaceObj = net.Interface{}
       interfaceSubstrings = []string{}
   }
   return interfaces
}

func (p *IOSXEParser) ParseInterfaceString(config string, interfaceType string) net.Interface {
   interfaceObj := net.Interface{ACL: make([]net.ACL, 0)}

   // Extract IPv4 address and subnet
   regexIPv4 := `\bip\s+address\s+((?:\d{1,3}\.){3}\d{1,3})\s+((?:\d{1,3}\.){3}\d{1,3})\b`
   ipv4Regex := regexp.MustCompile(regexIPv4)
   ipv4Match := ipv4Regex.FindStringSubmatch(config)
   ipv4Addr := ""
   ipv4Subnet := ""
   if len(ipv4Match) > 2 {
       ipv4Addr = ipv4Match[1]
       ipv4Subnet = ipv4Match[2]
   }
   // fmt.Printf("\n %s, %s", ipv4Addr, ipv4Subnet)

   // Check whether interface is in shutdown status
   regexShutdown := `\bshutdown\b`
   reShutdown := regexp.MustCompile(regexShutdown)
   reShutdownMatch := reShutdown.FindStringSubmatch(config)
   isShutdown := false
   if len(reShutdownMatch) >= 1 {
       isShutdown = true
   }

   // Extract Description
   regexDescription := regexp.MustCompile(`(?m)^\s*description\s+(.+)$`)
   reDescription := regexDescription.FindStringSubmatch(config)
   description := ""
   if len(reDescription) > 1 {
       description = reDescription[1]
   }

   // ACL
   acl := net.ACL{}
   regexAccessGroup := regexp.MustCompile(`ip access-group\s+(\S+)\s+(in|out)`)
   accessGroupMatch := regexAccessGroup.FindStringSubmatch(config)

   if accessGroupMatch != nil {
       acl.Name = accessGroupMatch[1]
       acl.Direction = accessGroupMatch[2]
       interfaceObj.ACL = append(interfaceObj.ACL, acl)
   } else {
       fmt.Println("No ip access-group match found")
   }

   interfaceObj.InterfaceType = interfaceType
   interfaceObj.Description = description
   interfaceObj.IP = ipv4Addr
   interfaceObj.Subnet = ipv4Subnet
   interfaceObj.IsShutdown = isShutdown

   return interfaceObj
}

func (p *IOSXEParser) ParseRouteMapBlock(blocks []string) []net.RouteMap {
   fmt.Println("Parsing Route-map block...")
   routeMaps := []net.RouteMap{}
   for _, block := range blocks {
       regexRouteMap := `route-map\s+([\w-]+)\s+(permit|deny)\s+(\d+)`
       reRouteMap := regexp.MustCompile(regexRouteMap)
       routeMapMatch := reRouteMap.FindStringSubmatch(block)
       // TODO
       // extract match and set rules
       // initially extract the prefix list in the "match ip address" statement
       // regexMatchIP := `^match\s+ip\s+address\s+prefix-list\s+(\S+)$`
       // reMatchIP := regexp.MustCompile(regexMatchIP)
       // matchIPMatches := reMatchIP.FindStringSubmatch(block)

       // routeMapMatch[1] - name of the route-map
       // routeMapMatch[2] - action (permit or deny)
       // routeMapMatch[3] - sequence number
       routeMapObj := net.RouteMap{
           Name:     routeMapMatch[1],
           Action:   routeMapMatch[2],
           Sequence: routeMapMatch[3],
       }
       routeMaps = append(routeMaps, routeMapObj)
   }
   return routeMaps
}

func (p *IOSXEParser) ParseIPAccessListBlock(blocks []string) []net.AccessList {
   fmt.Println("Parsing IP Access-list block...")
   accessLists := []net.AccessList{}
   for _, block := range blocks {
       regexIPAccessList := `ip\s+access-list\s+(standard|extended)\s+(\S+)`
       reAccessList := regexp.MustCompile(regexIPAccessList)
       accessListMatch := reAccessList.FindStringSubmatch(block)

       // extracted info
       aclType := accessListMatch[1]
       aclName := accessListMatch[2]

       accessListObj := net.AccessList{
           Name:  aclName,
           Type:  aclType,
           Rules: make([]net.ACLRule, 0),
       }

       // TODO
       // 1. "block" is a string. Try to iterate by newline instead of using finding all using regex
       lines := strings.Split(block, "\n")
       for j, line := range lines {
           if j == 0 {
               continue
           }
           var subnetMask string
           var dstIP string
           var action string
           if aclType == "standard" {
               aclRuleArr := strings.Split(strings.TrimSpace(line), " ")
               action = aclRuleArr[0]
               dstIP = aclRuleArr[1]
               if len(aclRuleArr) >= 3 {
                   subnetMask = aclRuleArr[2]
               } else {
                   subnetMask = ""
               }
               aclRuleObj := net.ACLRule{
                   Action:   action,
                   Type:     aclType,
                   SrcIP:    "",
                   SrcMask:  "",
                   DstIp:    dstIP,
                   DstMask:  subnetMask,
                   Protocol: "",
                   Port:     "",
               }
               accessListObj.Rules = append(accessListObj.Rules, aclRuleObj)
           } else {
               // TODO - Extended IP Access List
           }
       }
       accessLists = append(accessLists, accessListObj)
   }
   return accessLists
}

func (p *IOSXEParser) ParseIPPrefixListBlock(blocks []string) []net.PrefixList {
    fmt.Println("Parsing IP Prefix-list block...")
    prefixLists := []net.PrefixList{}
    for _, block := range blocks {
        // fmt.Printf("%v, %v \n", i, block)

        regexIPPrefixList := `ip prefix-list (\S+) seq (\d+) (permit|deny) (\d{1,3}(?:\.\d{1,3}){3})\/(\d{1,2})`
        rePrefixList := regexp.MustCompile(regexIPPrefixList)
        prefixListMatches := rePrefixList.FindAllStringSubmatch(block, -1)

        prefixListObj := net.PrefixList{
            Name:  "",
            Rules: make([]net.PrefixListRule, 0),
        }

        for i, match := range prefixListMatches {
            name := match[1]
            sequenceNumber := match[2]
            action := match[3]
            IPAddress := match[4]
            subnetMask := match[5]

            if i == 0 {
                prefixListObj.Name = name
            }

            prefixListRuleObj := net.PrefixListRule{
                SequenceNumber: sequenceNumber,
                Action: action,
                IP: IPAddress,
                Mask: subnetMask,
            }
            prefixListObj.Rules = append(prefixListObj.Rules, prefixListRuleObj)
        }

        prefixLists = append(prefixLists, prefixListObj)
    }
    return prefixLists
}

// TODO
func (p *IOSXEParser) ParseMetadata(s string) net.Metadata {
	// MetaData:
	// time source, hostname, version, current configuration bytes, last updated, last updated in NVRAM
	fmt.Println("Parse Metadata!")

	regexLastModified := `^! Last configuration change at (.+) by (\S+)`
	reLastModified := regexp.MustCompile(regexLastModified)

	regexHostname := `^hostname\s+(\S+)`
	reHostname := regexp.MustCompile(regexHostname)

	regexVersion := `^version\s+(\S+)`
	reVersion := regexp.MustCompile(regexVersion)

	scanner := bufio.NewScanner(strings.NewReader(s))
	currentLine := ""

	metadata := net.Metadata{}
	for scanner.Scan() {
		// prevLine := currentLine
		currentLine = scanner.Text()
		if matches := reLastModified.FindStringSubmatch(currentLine); matches != nil {
			metadata.LastModified = matches[1]
			metadata.LastModifiedUser = matches[2]
		} else if matches := reHostname.FindStringSubmatch(currentLine); matches != nil {
			metadata.Hostname = matches[1]
		} else if matches := reVersion.FindStringSubmatch(currentLine); matches != nil {
			metadata.Version = matches[1]
		}
	}

	return metadata
}