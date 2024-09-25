package net

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type BGP struct {
	ASN       string
	RouterID  string
	Neighbors []BGPNeighbor
}

type BGPNeighbor struct {
	ASN               string
	IP                string
	SessionAttributes map[string]string
}

func ParseBGPBlock(block string) BGP {
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

	bgpObject := BGP{
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
			bgpNeighborObject := BGPNeighbor{
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
