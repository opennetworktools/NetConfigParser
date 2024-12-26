package net

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
