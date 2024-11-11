package net

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

type PrefixList struct {
    Name string
    Rules []PrefixListRule
}

type PrefixListRule struct {
    SequenceNumber string
    Action string
    IP     string
    Mask   string
}

