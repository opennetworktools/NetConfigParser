package net

type Interface struct {
	InterfaceType string
	IP            string
	Description   string
}

type Interfaces struct {
	GigabitEthernet []Interface
	PortChannel     []Interface
	Loopback        []Interface
	Tunnel          []Interface
}
