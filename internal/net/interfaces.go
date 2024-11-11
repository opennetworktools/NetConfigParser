package net

type Interface struct {
	InterfaceType string
	Description   string
	IP            string
	Subnet        string
	IsShutdown    bool
	ACL           []ACL
}

type Interfaces struct {
	GigabitEthernet    []Interface
	TenGigabitEthernet []Interface
	PortChannel        []Interface
	Loopback           []Interface
	Tunnel             []Interface
}

type ACL struct {
	Name      string
	Direction string
}
