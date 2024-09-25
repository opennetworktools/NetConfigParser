# NetConfigParser

NetConfigParser is a tool that parses your network device's startup or running configurations.

## Design

### TODO:

1. Extract all the features from the running configs as an array of strings
2. Implement parsing logic for,
- AAA
- Interfaces (Ethernet, PortChannel, Loopback etc.)
- OSPF
- BGP
- Prefix-list (both IPv4 and IPv6), Access-list, Route-map etc.