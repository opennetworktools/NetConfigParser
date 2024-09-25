package net

import (
	"fmt"
	"regexp"
)

type RouteMap struct {
	Name     string
	Action   string
	Sequence string
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
