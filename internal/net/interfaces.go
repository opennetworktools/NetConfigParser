package net

import (
	"fmt"
	"regexp"
)

type Interface struct {
	InterfaceType string
	Description   string
	IP            string
	Subnet        string
	IsShutdown    bool
	ACL 		  []ACL
}

type Interfaces struct {
	GigabitEthernet    []Interface
	TenGigabitEthernet []Interface
	PortChannel        []Interface
	Loopback           []Interface
	Tunnel             []Interface
}

type ACL struct {
	Name 	  string
	Direction string
}

func ParseInterfacesBlock(blocks []string) Interfaces {
	fmt.Println("Parsing Interfaces Block...")

	interfaces := Interfaces{
		GigabitEthernet:    []Interface{},
		TenGigabitEthernet: []Interface{},
		PortChannel:        []Interface{},
		Loopback:           []Interface{},
		Tunnel:             []Interface{},
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

	var interfaceObj Interface
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
			interfaceObj = parseInterfaceString(block, interfaceSubstrings[0])
			interfaces.GigabitEthernet = append(interfaces.GigabitEthernet, interfaceObj)
			// fmt.Println(interfaceObj)
		} else if reTenGigabitEthernet.MatchString(block) {
			fmt.Println("Interface TenGigabitEthernet Matched!")
			interfaceSubstrings = reTenGigabitEthernet.FindStringSubmatch(block)
			interfaceObj = parseInterfaceString(block, interfaceSubstrings[0])
			interfaces.TenGigabitEthernet = append(interfaces.TenGigabitEthernet, interfaceObj)
		} else if reLoopback.MatchString(block) {
			fmt.Println("Interface Loopback Matched!")
			interfaceSubstrings = reLoopback.FindStringSubmatch(block)
			interfaceObj = parseInterfaceString(block, interfaceSubstrings[0])
			interfaces.Loopback = append(interfaces.Loopback, interfaceObj)
		} else if rePortChannel.MatchString(block) {
			fmt.Println("Interface PortChannel Matched!")
			interfaceSubstrings = rePortChannel.FindStringSubmatch(block)
			interfaceObj = parseInterfaceString(block, interfaceSubstrings[0])
			interfaces.PortChannel = append(interfaces.PortChannel, interfaceObj)
		} else if reTunnel.MatchString(block) {
			fmt.Println("Interface Tunnel Matched!")
			interfaceSubstrings = reTunnel.FindStringSubmatch(block)
			interfaceObj = parseInterfaceString(block, interfaceSubstrings[0])
			interfaces.Tunnel = append(interfaces.Tunnel, interfaceObj)
		}
		interfaceObj = Interface{}
		interfaceSubstrings = []string{}
	}
	return interfaces
}

func parseInterfaceString(config string, interfaceType string) Interface {
	interfaceObj := Interface{ACL: make([]ACL, 0)}
	
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
	acl := ACL{}
	regexAccessGroup := regexp.MustCompile(`ip access-group\s+(\S+)\s+(in|out)`)
	accessGroupMatch := regexAccessGroup.FindStringSubmatch(config)

	if accessGroupMatch != nil {
		acl.Name = accessGroupMatch[1]
		acl.Direction = accessGroupMatch[2]
		interfaceObj.ACL = append(interfaceObj.ACL, acl)
	} else {
		// fmt.Println("No ip access-group match found")
	}

	interfaceObj.InterfaceType = interfaceType
	interfaceObj.Description = description
	interfaceObj.IP = ipv4Addr
	interfaceObj.Subnet = ipv4Subnet
	interfaceObj.IsShutdown = isShutdown

	return interfaceObj
}
