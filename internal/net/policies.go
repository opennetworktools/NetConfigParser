package net

import (
	"fmt"
	"regexp"
	"strings"
)

type RouteMap struct {
	Name     string
	Action   string
	Sequence string
}

type AccessList struct {
    Name  string
    Type  string
    Rules []ACLRule
}

type ACLRule struct {
    Action   string
    Type     string
    SrcIP    string
    SrcMask  string
    DstIp    string
    DstMask  string
    Protocol string
    Port     string
}

func ParseRouteMapBlock(blocks []string) []RouteMap {
	fmt.Println("Parsing Route-map block...")
	routeMaps := []RouteMap{}
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
		routeMapObj := RouteMap{
			Name:     routeMapMatch[1],
			Action:   routeMapMatch[2],
			Sequence: routeMapMatch[3],
		}
		routeMaps = append(routeMaps, routeMapObj)
	}
	return routeMaps
}

func ParseIPAccessListBlock(blocks []string) []AccessList {
    fmt.Println("Parsing IP Access-list block...")
    accessLists := []AccessList{}
    for _, block := range blocks {
        regexIPAccessList := `ip\s+access-list\s+(standard|extended)\s+(\S+)`
        reAccessList := regexp.MustCompile(regexIPAccessList)
        accessListMatch := reAccessList.FindStringSubmatch(block)

        // extracted info
        aclType := accessListMatch[1]
        aclName := accessListMatch[2]

        accessListObj := AccessList{
            Name:  aclName,
            Type:  aclType,
            Rules: make([]ACLRule, 0),
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
                aclRuleObj := ACLRule{
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
